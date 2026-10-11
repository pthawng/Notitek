"""Validate the proposed API artifacts; this does not test a running service.

python -m pip install -r scripts/requirements-contract.txt
python scripts/validate_notification_contract.py
"""
import copy
import json
from pathlib import Path

from jsonschema import Draft202012Validator, FormatChecker, ValidationError
from openapi_spec_validator import validate

ROOT = Path(__file__).resolve().parents[1]
API = ROOT / 'docs' / 'api'
SPEC = json.loads((API / 'openapi.json').read_text(encoding='utf-8'))


def validator(name):
    return Draft202012Validator(
        {'$ref': '#/components/schemas/' + name, 'components': SPEC['components']},
        format_checker=FormatChecker(),
    )


def check_request(body):
    validator('SendRequest').validate(body)
    profiles = SPEC['components']['schemas']['SendRequest']['x-type-profiles']
    profile = profiles.get(f"{body['type']}:{body['type_version']}")
    if profile is None:
        raise ValidationError('Unregistered type/version')
    validator(profile['$ref'].split('/')[-1]).validate(body)


def main():
    validate(SPEC)
    assert set(SPEC['paths']) == {
        '/v1/notifications', '/v1/notifications/{notification_id}',
        '/v1/notifications/cancel', '/v1/notifications/{notification_id}/cancel',
    }
    order_types = {'created', 'cancelled', 'pickup_failed', 'picked_up', 'delivery_failed', 'delivered',
                   'partially_delivered', 'return_confirmed', 'returned', 'sla_breached'}
    profiles = SPEC['components']['schemas']['SendRequest']['x-type-profiles']
    assert {f'order.{t}:1' for t in order_types} <= set(profiles), 'All 10 Order milestones must be registered'
    requests = {}
    for path in sorted((API / 'examples').glob('*.request.json')):
        body = json.loads(path.read_text(encoding='utf-8'))
        check_request(body)
        requests[path.name.split('.')[0]] = body
    for path in sorted((API / 'examples').glob('*.response.json')):
        body = json.loads(path.read_text(encoding='utf-8'))
        validator('ErrorResponse' if body['error'] else 'NotificationResponse').validate(body)

    cases = []

    def bad(name, example, change):
        body = copy.deepcopy(requests[example])
        change(body)
        cases.append((name, body))

    bad('OTP cannot queue', 'otp', lambda b: b.update(mode='async'))
    bad('OTP needs secret', 'otp', lambda b: b.pop('secrets'))
    bad('OTP cannot include arbitrary secret', 'otp', lambda b: b['secrets'].update(token='x'))
    bad('OTP phone must be E.164', 'otp', lambda b: b['recipients'][0].update(phone='0901234567'))
    bad('ACCOUNT cannot carry Shop', 'otp', lambda b: b['context'].update(organization_id='org-x'))
    bad('Source cannot be spoofed in body', 'otp', lambda b: b.update(source='order'))
    bad('No raw channel override', 'otp', lambda b: b.update(channels=['SMS']))
    bad('Unknown type', 'otp', lambda b: b.update(type='unregistered.type'))
    bad('Unknown type version', 'otp', lambda b: b.update(type_version=2))
    bad('Offset required', 'otp', lambda b: b.update(occurred_at='2026-10-11T10:00:00'))
    bad('Missing deadline', 'otp', lambda b: b.pop('send_deadline'))
    bad('Too many recipients', 'otp', lambda b: b.update(recipients=b['recipients'] * 101))
    bad('Invitation must carry organization', 'invitation', lambda b: b['context'].pop('organization_id'))
    bad('Legacy shop_id is rejected', 'invitation', lambda b: b['context'].update(shop_id='shop-x'))
    bad('camelCase fields are rejected', 'invitation', lambda b: b.update(sourceRef=b.pop('source_ref')))
    bad('Email must be email', 'invitation', lambda b: b['recipients'][0].update(email='invalid'))
    bad('Link requires HTTPS', 'invitation', lambda b: b['secrets'].update(invitation_url='http://example.invalid'))
    bad('No extra template fields', 'invitation', lambda b: b['data'].update(html='<b>x</b>'))
    bad('No arbitrary Shop audience', 'order-delivery-failed', lambda b: b['recipients'][0].update(audience_code='ALL_USERS'))
    bad('Order needs consignee contact', 'order-delivery-failed', lambda b: b['recipients'].pop())
    bad('Order cannot replace audience with contact', 'order-delivery-failed', lambda b: b['recipients'].__setitem__(0, copy.deepcopy(b['recipients'][1])))
    bad('Order requires cancellation group', 'order-delivery-failed', lambda b: b.pop('cancel_group'))
    bad('Order needs Shop audience', 'order-created', lambda b: b['recipients'].pop(0))
    bad('Order consignee needs phone', 'order-delivery-failed', lambda b: b['recipients'][1].pop('phone'))
    bad('Order status code format', 'order-created', lambda b: b['data'].update(order_status='0301'))
    bad('Order tracking link requires HTTPS', 'order-created', lambda b: b['data'].update(tracking_url='http://example.invalid'))
    bad('Order cannot carry secret', 'order-created', lambda b: b.update(secrets={'otp':'123456'}))
    bad('Order milestone without failure cannot carry cancel group', 'order-delivered', lambda b: b.update(cancel_group='x'))
    bad('Only one ops audience', 'order-sla-breached', lambda b: b['recipients'].append(copy.deepcopy(b['recipients'][-1])))
    bad('SLA breach needs expected time', 'order-sla-breached', lambda b: b['data'].pop('expected_to'))
    bad('Pickup failure needs failure id', 'order-pickup-failed', lambda b: b['data'].pop('pickup_failure_id'))
    bad('Order needs organization', 'order-returned', lambda b: b['context'].pop('organization_id'))
    bad('Order needs subject_ref', 'order-created', lambda b: b.pop('subject_ref'))
    bad('subject_ref has no PII format', 'order-created', lambda b: b.update(subject_ref='+84 900 000 000'))
    bad('Device alert requires account', 'device-alert', lambda b: b['recipients'].__setitem__(0, {'kind':'CONTACT','recipient_ref':'x','phone':'+84900000000'}))
    bad('Device alert cannot carry secret', 'device-alert', lambda b: b.update(secrets={'otp':'123456'}))
    bad('OTP needs purpose', 'otp', lambda b: b['data'].pop('purpose'))
    bad('OTP purpose is enumerated', 'otp', lambda b: b['data'].update(purpose='MARKETING'))
    bad('OTP email must be email', 'otp-email', lambda b: b['recipients'][0].update(email='invalid'))
    bad('OTP cannot target account id', 'otp', lambda b: b['recipients'].__setitem__(0, {'kind':'ACCOUNT','recipient_ref':'x','account_id':'acc-x'}))
    bad('Member invitation needs inviter', 'invitation', lambda b: b['data'].pop('inviter_name'))
    bad('Invitation kind is enumerated', 'invitation', lambda b: b['data'].update(invitation_kind='GUEST'))
    bad('Activation link requires HTTPS', 'account-activation', lambda b: b['secrets'].update(activation_url='http://example.invalid'))
    bad('Activation needs expiry', 'account-activation', lambda b: b.pop('expires_at'))
    bad('Activation is email only', 'account-activation', lambda b: b['recipients'].__setitem__(0, {'kind':'CONTACT','recipient_ref':'x','phone':'+84900000000'}))
    bad('Password notice has at most two contacts', 'password-changed', lambda b: b['recipients'].append(copy.deepcopy(b['recipients'][0])))
    bad('Password notice cannot carry secret', 'password-changed', lambda b: b.update(secrets={'otp':'123456'}))
    bad('Password notice changed_via is enumerated', 'password-changed', lambda b: b['data'].update(changed_via='ADMIN'))
    bad('Identifier notice needs masked value', 'identifier-changed', lambda b: b['data'].pop('masked_new_value'))
    bad('Offboarding needs organization scope', 'employee-offboarded', lambda b: b['context'].pop('organization_id'))
    bad('Offboarding effective_at needs offset', 'employee-offboarded', lambda b: b['data'].update(effective_at='2026-10-11T17:00:00'))

    for name, body in cases:
        try:
            check_request(body)
        except ValidationError:
            continue
        raise AssertionError('Invalid request was accepted: ' + name)
    page = {'error': False, 'message': 'OK', 'data': {
        'items': [], 'meta': {'page_size': 20, 'has_next': False, 'next_cursor': None}}}
    validator('NotificationPageResponse').validate(page)
    delivery = copy.deepcopy(json.loads((API / 'examples' / 'otp-accepted.response.json').read_text(encoding='utf-8'))['data'])
    if delivery['deliveries']:
        delivery['deliveries'][0].pop('attempt_count')
        try:
            validator('Notification').validate(delivery)
            raise AssertionError('Delivery without attempt_count was accepted')
        except ValidationError:
            pass
    validator('CancelRequest').validate({'reason_code': 'SOURCE_RESOLVED'})
    validator('CancelGroupRequest').validate({
        'reason_code': 'SOURCE_RESOLVED', 'cancel_group': 'failure-demo',
        'context': requests['order-delivery-failed']['context'],
    })
    print(f'PASS: OpenAPI, {len(requests)} request examples, 4 response examples, {len(cases)} negative cases, 2 cancel bodies; canonical route /v1/notifications')
    print('Runtime authorization, timing, UTF-8 byte limits, concurrency and providers are not tested.')


if __name__ == '__main__':
    main()
