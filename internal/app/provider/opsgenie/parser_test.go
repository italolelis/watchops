package opsgenie_test

import (
	"testing"
	"time"

	"github.com/italolelis/watchops/internal/app/provider"
	"github.com/italolelis/watchops/internal/app/provider/opsgenie"
	"github.com/stretchr/testify/assert"
)

func TestParsePayload(t *testing.T) {
	cases := []struct {
		name           string
		payload        []byte
		expectedObject provider.Event
		shouldFail     bool
	}{
		{
			name:    "empty header",
			payload: []byte(`{"action":"AddNote","alert":{"alertId":"f3efa74b-29cc-4a4f-9bb7-804d7c1be544-1639476708514","message":"[Prometheus]: [FIRING:1]  (NonCompletedDeletions lykon-service-0.8.3 app http-probe 10.0.3.129:9090 gdpr-deleter-srv gdpr gdpr-del","tags":["service=gdpr-deleter-srv"],"tinyId":"444","entity":"","alias":"94d30522a04f65515397ec16b26d4e0fee5946671e0c985ac07aca345aba1968","createdAt":1639476708514,"updatedAt":1639521280284000000,"username":"italo.vietro@lykon.com","userId":"5e7ececc-df4a-40f1-9c6c-3598c90422b0","note":"Test 2","description":"There are 1 non-completed deletions in crm-srv for more than 10m https://github.com/vimeda/runbooks Service did not complete requested deletions in 24h\nAlerts Firing:\nLabels:\n - alertname = NonCompletedDeletions\n - chart = lykon-service-0.8.3\n - container = app\n - endpoint = http-probe\n - instance = 10.0.3.129:9090\n - job = gdpr-deleter-srv\n - namespace = gdpr\n - pod = gdpr-deleter-srv-app-965459dc4-x9dnd\n - prometheus = monitoring/kube-prometheus-kube-prome-prometheus\n - release = gdpr-deleter-srv\n - service = gdpr-deleter-srv\n - service_name = crm-srv\n - severity = page\nAnnotations:\n - description = There are 1 non-completed deletions in crm-srv for more than 10m\n - runbook_url = https://github.com/vimeda/runbooks\n - summary = Service did not complete requested deletions in 24h\nSource: http://kube-prometheus-kube-prome-prometheus.monitoring:9090/graph?g0.expr=gdpr_deleter_srv_lykon_io_deleter_checker+%3E+0&g0.tab=1","responders":[{"id":"a6da4929-aa71-41bc-ac30-ac7ae9519034","type":"team","name":"Tech"}],"teams":["a6da4929-aa71-41bc-ac30-ac7ae9519034"],"actions":[],"details":{"alertname":"NonCompletedDeletions","chart":"lykon-service-0.8.3","container":"app","endpoint":"http-probe","instance":"10.0.3.129:9090","job":"gdpr-deleter-srv","namespace":"gdpr","pod":"gdpr-deleter-srv-app-965459dc4-x9dnd","prometheus":"monitoring/kube-prometheus-kube-prome-prometheus","release":"gdpr-deleter-srv","service":"gdpr-deleter-srv","service_name":"crm-srv","severity":"page"},"priority":"P1","oldPriority":"P1","source":"http://kube-prometheus-kube-prome-alertmanager.monitoring:9093/#/alerts?receiver=opsgenie"},"source":{"name":"","type":"web"},"integrationName":"WatchOps Webhook","integrationId":"04244f46-c6d3-4069-9f7f-cd3c5cc63120","integrationType":"Webhook"}`),
			expectedObject: provider.Event{
				ID:          "f3efa74b-29cc-4a4f-9bb7-804d7c1be544-1639476708514",
				EventType:   "AddNote",
				TimeCreated: time.UnixMicro(1639521280284000000),
				Signature:   "l1Q61kXRlCp_qTIudME8AfYmVtY=",
				MsgID:       "testID",
				Source:      "opsgenie",
				Metadata:    []byte(`{"action":"AddNote","alert":{"alertId":"f3efa74b-29cc-4a4f-9bb7-804d7c1be544-1639476708514","message":"[Prometheus]: [FIRING:1]  (NonCompletedDeletions lykon-service-0.8.3 app http-probe 10.0.3.129:9090 gdpr-deleter-srv gdpr gdpr-del","tags":["service=gdpr-deleter-srv"],"tinyId":"444","entity":"","alias":"94d30522a04f65515397ec16b26d4e0fee5946671e0c985ac07aca345aba1968","createdAt":1639476708514,"updatedAt":1639521280284000000,"username":"italo.vietro@lykon.com","userId":"5e7ececc-df4a-40f1-9c6c-3598c90422b0","note":"Test 2","description":"There are 1 non-completed deletions in crm-srv for more than 10m https://github.com/vimeda/runbooks Service did not complete requested deletions in 24h\nAlerts Firing:\nLabels:\n - alertname = NonCompletedDeletions\n - chart = lykon-service-0.8.3\n - container = app\n - endpoint = http-probe\n - instance = 10.0.3.129:9090\n - job = gdpr-deleter-srv\n - namespace = gdpr\n - pod = gdpr-deleter-srv-app-965459dc4-x9dnd\n - prometheus = monitoring/kube-prometheus-kube-prome-prometheus\n - release = gdpr-deleter-srv\n - service = gdpr-deleter-srv\n - service_name = crm-srv\n - severity = page\nAnnotations:\n - description = There are 1 non-completed deletions in crm-srv for more than 10m\n - runbook_url = https://github.com/vimeda/runbooks\n - summary = Service did not complete requested deletions in 24h\nSource: http://kube-prometheus-kube-prome-prometheus.monitoring:9090/graph?g0.expr=gdpr_deleter_srv_lykon_io_deleter_checker+%3E+0&g0.tab=1","responders":[{"id":"a6da4929-aa71-41bc-ac30-ac7ae9519034","type":"team","name":"Tech"}],"teams":["a6da4929-aa71-41bc-ac30-ac7ae9519034"],"actions":[],"details":{"alertname":"NonCompletedDeletions","chart":"lykon-service-0.8.3","container":"app","endpoint":"http-probe","instance":"10.0.3.129:9090","job":"gdpr-deleter-srv","namespace":"gdpr","pod":"gdpr-deleter-srv-app-965459dc4-x9dnd","prometheus":"monitoring/kube-prometheus-kube-prome-prometheus","release":"gdpr-deleter-srv","service":"gdpr-deleter-srv","service_name":"crm-srv","severity":"page"},"priority":"P1","oldPriority":"P1","source":"http://kube-prometheus-kube-prome-alertmanager.monitoring:9093/#/alerts?receiver=opsgenie"},"source":{"name":"","type":"web"},"integrationName":"WatchOps Webhook","integrationId":"04244f46-c6d3-4069-9f7f-cd3c5cc63120","integrationType":"Webhook"}`),
			},
			shouldFail: false,
		},
	}

	p := &opsgenie.Parser{}
	headers := make(map[string][]string)
	headers["msg_id"] = []string{"testID"}

	for _, c := range cases {
		c := c

		t.Run(c.name, func(t *testing.T) {
			e, err := p.Parse(headers, c.payload)
			if c.shouldFail {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, c.expectedObject, e)
			}
		})
	}
}
