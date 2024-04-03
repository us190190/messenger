
COMPUTE_URL_BASE = 'https://www.googleapis.com/compute/v1/'

def GenerateConfig(context):

  resources = [{
    'name': context.properties['servername'],
    'type': 'compute.v1.instance',
    'properties': {
      'zone': context.properties['zone'],
      'machineType': ''.join([COMPUTE_URL_BASE, 'projects/', context.properties['project'],
                              '/zones/', context.properties['zone'], '/',
                              'machineTypes/e2-small']),
      'disks': [{
        'deviceName': 'boot',
        'type': 'PERSISTENT',
        'boot': True,
        'sizeGb': 25,
        'autoDelete': True,
        'initializeParams': {
          'sourceImage': ''.join([COMPUTE_URL_BASE, 'projects/',
                                  'debian-cloud/global',
                                  '/images/family/debian-11'])
        }
      }],
      'metadata': {
        "items": [
          {
            "key": "startup-script",
            "value": context.properties['startup-script']
          }
        ]
      },
      'tags': context.properties['tags'],
      'networkInterfaces': [{
        'network': ''.join([COMPUTE_URL_BASE, 'projects/',
                            context.env['project'],
                            '/global/networks/', context.properties['network']]),
        'accessConfigs': [{
          'name': 'External NAT',
          'type': 'ONE_TO_ONE_NAT'
        }]
      }],
      'labels': {
        "name": context.properties['servername'],
        "environment": context.properties['environment']
      }
    }
  }]
  return {'resources': resources}