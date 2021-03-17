{
  grafanaDashboards: {
    'couchdb-overview.json': (import 'dashboards/couchdb-overview.json'),
  },

  prometheusAlerts+:
    importRules(importstr 'alerts/general.yaml')
}