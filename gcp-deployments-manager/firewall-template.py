"""Creates the firewall."""


def GenerateConfig(context):
    """Creates the firewall."""

    resources = [{
        'name': context.env['name'],
        'type': 'compute.v1.firewall',
        'properties': {
            'network': 'global/networks/default',
            'targetTags': ["https"],
            'sourceRanges': ['0.0.0.0/0'],
            'allowed': [{
                'IPProtocol': 'TCP',
                'ports': [443]
            }]
        }
    }]
    return {'resources': resources}