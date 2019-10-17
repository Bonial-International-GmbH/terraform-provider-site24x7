// TODO
// This is work in progress and have to be completed

package site24x7

//func TestWebsiteMonitorDelete(t *testing.T) {
//	d := monitorTestResourceData(t)
//	d.SetId("123")
//
//	c := fake.NewClient()
//
//	c.FakeMonitors.On("Delete", "123").Return(nil).Once()
//
//	require.NoError(t, monitorDelete(d, c))
//}

//func monitorTestResourceData(t *testing.T) *schema.ResourceData {
//	//return schema.TestResourceDataRaw(t, MonitorSchema, map[string]interface{}{
//	return schema.TestResourceDataRaw(t, nil, map[string]interface{}{
//		"DisplayNmae":    "foo",
//		"Type":           "URL",
//		"Website":        "www.test.tld",
//		"CheckFrequency": "60",
//		"HTTPMethod":     "P",
//		"AuthUser":       "username",
//		"AuthPass":       "password",
//		"MatchCase":      true,
//		"UserAgent":      "firefox",
//		"CustomHeaders": map[string]string{
//			"Name":  "Accept-Encoding",
//			"Value": "gzip",
//		},
//		"Timeout":               120,
//		"LocationProfileID":     "456",
//		"NotificationProfileID": "789",
//		"ThresholdProfileID":    "012",
//		"MonitorGroups": []string{
//			"234",
//			"567",
//		},
//		"UserGroupIDs": []string{
//			"123",
//			"456",
//		},
//		"UseNameServer": true,
//	})
//}
