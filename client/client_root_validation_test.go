package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/docker/notary/keystoremanager"
	"github.com/endophage/gotuf/data"
	"github.com/stretchr/testify/assert"
)

type SignedRSARootTemplate struct {
	RootPem string
}

var passphraseRetriever = func (string, string, bool, int) (string, bool, error) { return "passphrase", false, nil}

const validPEMEncodedRSARoot = `LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUZLekNDQXhXZ0F3SUJBZ0lRUnlwOVFxY0pmZDNheXFkaml6OHhJREFMQmdrcWhraUc5dzBCQVFzd09ERWEKTUJnR0ExVUVDaE1SWkc5amEyVnlMbU52YlM5dWIzUmhjbmt4R2pBWUJnTlZCQU1URVdSdlkydGxjaTVqYjIwdgpibTkwWVhKNU1CNFhEVEUxTURjeE56QTJNelF5TTFvWERURTNNRGN4TmpBMk16UXlNMW93T0RFYU1CZ0dBMVVFCkNoTVJaRzlqYTJWeUxtTnZiUzl1YjNSaGNua3hHakFZQmdOVkJBTVRFV1J2WTJ0bGNpNWpiMjB2Ym05MFlYSjUKTUlJQ0lqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FnOEFNSUlDQ2dLQ0FnRUFvUWZmcnpzWW5zSDh2R2Y0Smg1NQpDajV3cmpVR3pEL3NIa2FGSHB0ako2VG9KR0p2NXlNQVB4enlJbnU1c0lvR0xKYXBuWVZCb0FVMFlnSTlxbEFjCllBNlN4YVN3Z202cnB2bW5sOFFuMHFjNmdlcjNpbnBHYVVKeWxXSHVQd1drdmNpbVFBcUhaeDJkUXRMN2c2a3AKcm1LZVRXcFdvV0x3M0pvQVVaVVZoWk1kNmEyMlpML0R2QXcrSHJvZ2J6NFhleWFoRmI5SUg0MDJ6UHhONnZnYQpKRUZURjBKaTFqdE5nME1vNHBiOVNIc01zaXcrTFpLN1NmZkhWS1B4dmQyMW0vYmlObXdzZ0V4QTNVOE9PRzhwCnV5Z2ZhY3lzNWM4K1pyWCtaRkcvY3Z3S3owazYvUWZKVTQwczZNaFh3NUMyV3R0ZFZtc0c5LzdyR0ZZakhvSUoKd2VEeXhnV2s3dnhLelJKSS91bjdjYWdESWFRc0tySlFjQ0hJR0ZSbHBJUjVUd1g3dmwzUjdjUm5jckRSTVZ2YwpWU0VHMmVzeGJ3N2p0eklwL3lwblZSeGNPbnk3SXlweWpLcVZlcVo2SGd4WnRUQlZyRjFPL2FIbzJrdmx3eVJTCkF1czRrdmg2ejMranpUbTlFemZYaVBRelk5QkVrNWdPTHhoVzlyYzZVaGxTK3BlNWxrYU4vSHlxeS9sUHVxODkKZk1yMnJyN2xmNVdGZEZuemU2V05ZTUFhVzdkTkE0TkUwZHlENTM0MjhaTFh4TlZQTDRXVTY2R2FjNmx5blE4bApyNXRQc1lJRlh6aDZGVmFSS0dRVXRXMWh6OWVjTzZZMjdSaDJKc3lpSXhnVXFrMm9veEU2OXVONDJ0K2R0cUtDCjFzOEcvN1Z0WThHREFMRkxZVG56THZzQ0F3RUFBYU0xTURNd0RnWURWUjBQQVFIL0JBUURBZ0NnTUJNR0ExVWQKSlFRTU1Bb0dDQ3NHQVFVRkJ3TURNQXdHQTFVZEV3RUIvd1FDTUFBd0N3WUpLb1pJaHZjTkFRRUxBNElDQVFCTQpPbGwzRy9YQno4aWRpTmROSkRXVWgrNXczb2ptd2FuclRCZENkcUVrMVdlbmFSNkR0Y2ZsSng2WjNmL213VjRvCmIxc2tPQVgxeVg1UkNhaEpIVU14TWljei9RMzhwT1ZlbEdQclduYzNUSkIrVktqR3lIWGxRRFZrWkZiKzQrZWYKd3RqN0huZ1hoSEZGRFNnam0zRWRNbmR2Z0RRN1NRYjRza09uQ05TOWl5WDdlWHhoRkJDWm1aTCtIQUxLQmoyQgp5aFY0SWNCRHFtcDUwNHQxNHJ4OS9KdnR5MGRHN2ZZN0k1MWdFUXBtNFMwMkpNTDV4dlRtMXhmYm9XSWhaT0RJCnN3RUFPK2VrQm9GSGJTMVE5S01QaklBdzNUckNISDh4OFhacTV6c1l0QUMxeVpIZENLYTI2YVdkeTU2QTllSGoKTzFWeHp3bWJOeVhSZW5WdUJZUCswd3IzSFZLRkc0Sko0WlpwTlp6UVcvcHFFUGdoQ1RKSXZJdWVLNjUyQnlVYwovL3N2K25YZDVmMTlMZUVTOXBmMGwyNTNORGFGWlBiNmFlZ0tmcXVXaDhxbFFCbVVRMkd6YVRMYnRtTmQyOE02Clc3aUw3dGtLWmUxWm5CejlSS2d0UHJEampXR1pJbmpqY09VOEV0VDRTTHE3a0NWRG1QczVNRDh2YUFtOTZKc0UKam1MQzNVdS80azdIaURZWDBpMG1PV2tGalpRTWRWYXRjSUY1RlBTcHB3c1NiVzhRaWRuWHQ1NFV0d3RGREVQegpscGpzN3liZVFFNzFKWGNNWm5WSUs0YmpSWHNFRlBJOThScElsRWRlZGJTVWRZQW5jTE5KUlQ3SFpCTVBHU3daCjBQTkp1Z2xubHIzc3JWemRXMWR6MnhRamR2THd4eTZtTlVGNnJiUUJXQT09Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K`

const validCAPEMEncodeRSARoot = `LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tDQpNSUlHTXpDQ0JCdWdBd0lCQWdJQkFUQU5CZ2txaGtpRzl3MEJBUXNGQURCZk1Rc3dDUVlEVlFRR0V3SlZVekVMDQpNQWtHQTFVRUNBd0NRMEV4RmpBVUJnTlZCQWNNRFZOaGJpQkdjbUZ1WTJselkyOHhEekFOQmdOVkJBb01Ca1J2DQpZMnRsY2pFYU1CZ0dBMVVFQXd3UlRtOTBZWEo1SUZSbGMzUnBibWNnUTBFd0hoY05NVFV3TnpFMk1EUXlOVEF6DQpXaGNOTWpVd056RXpNRFF5TlRBeldqQmZNUm93R0FZRFZRUUREQkZPYjNSaGNua2dWR1Z6ZEdsdVp5QkRRVEVMDQpNQWtHQTFVRUJoTUNWVk14RmpBVUJnTlZCQWNNRFZOaGJpQkdjbUZ1WTJselkyOHhEekFOQmdOVkJBb01Ca1J2DQpZMnRsY2pFTE1Ba0dBMVVFQ0F3Q1EwRXdnZ0lpTUEwR0NTcUdTSWIzRFFFQkFRVUFBNElDRHdBd2dnSUtBb0lDDQpBUUN3VlZENHBLN3o3cFhQcEpiYVoxSGc1ZVJYSWNhWXRiRlBDbk4waXF5OUhzVkVHbkVuNUJQTlNFc3VQK20wDQo1TjBxVlY3REdiMVNqaWxvTFhEMXFERHZoWFdrK2dpUzlwcHFQSFBMVlBCNGJ2enNxd0RZcnRwYnFrWXZPMFlLDQowU0wza3hQWFVGZGxrRmZndTB4amxjem0yUGhXRzNKZDhhQXRzcEwvTCtWZlBBMTNKVWFXeFNMcHVpMUluOHJoDQpnQXlRVEs2UTRPZjZHYkpZVG5BSGI1OVVvTFhTekI1QWZxaVVxNkw3bkVZWUtvUGZsUGJSQUlXTC9VQm0wYytIDQpvY21zNzA2UFlwbVBTMlJRdjNpT0dtbm45aEVWcDNQNmpxN1dBZXZiQTRhWUd4NUVzYlZ0WUFCcUpCYkZXQXV3DQp3VEdSWW16bjBNajBlVE1nZTl6dFlCMi8yc3hkVGU2dWhtRmdwVVhuZ0RxSkk1TzlOM3pQZnZsRUltQ2t5M0hNDQpqSm9MN2c1c21xWDlvMVArRVNMaDBWWnpoaDdJRFB6UVRYcGNQSVMvNnowbDIyUUdrSy8xTjFQYUFEYVVIZExMDQp2U2F2M3kyQmFFbVB2ZjJma1pqOHlQNWVZZ2k3Q3c1T05oSExEWUhGY2w5Wm0veXdtZHhISkVUejluZmdYbnNXDQpITnhEcXJrQ1ZPNDZyL3U2clNyVXQ2aHIzb2RkSkc4czhKbzA2ZWFydzZYVTNNek0rM2dpd2tLMFNTTTN1UlBxDQo0QXNjUjFUditFMzFBdU9BbWpxWVFvVDI5Yk1JeG9TemVsamovWW5lZHdqVzQ1cFd5YzNKb0hhaWJEd3ZXOVVvDQpHU1pCVnk0aHJNL0ZhN1hDV3YxV2ZITlcxZ0R3YUxZd0RubDVqRm1SQnZjZnVRSURBUUFCbzRINU1JSDJNSUdSDQpCZ05WSFNNRWdZa3dnWWFBRkhVTTFVM0U0V3lMMW52RmQrZFBZOGY0TzJoWm9XT2tZVEJmTVFzd0NRWURWUVFHDQpFd0pWVXpFTE1Ba0dBMVVFQ0F3Q1EwRXhGakFVQmdOVkJBY01EVk5oYmlCR2NtRnVZMmx6WTI4eER6QU5CZ05WDQpCQW9NQmtSdlkydGxjakVhTUJnR0ExVUVBd3dSVG05MFlYSjVJRlJsYzNScGJtY2dRMEdDQ1FEQ2VETGJlbUlUDQpTekFTQmdOVkhSTUJBZjhFQ0RBR0FRSC9BZ0VBTUIwR0ExVWRKUVFXTUJRR0NDc0dBUVVGQndNQ0JnZ3JCZ0VGDQpCUWNEQVRBT0JnTlZIUThCQWY4RUJBTUNBVVl3SFFZRFZSME9CQllFRkhlNDhoY0JjQXAwYlVWbFR4WGVSQTRvDQpFMTZwTUEwR0NTcUdTSWIzRFFFQkN3VUFBNElDQVFBV1V0QVBkVUZwd1JxK04xU3pHVWVqU2lrZU1HeVBac2NaDQpKQlVDbWhab0Z1ZmdYR2JMTzVPcGNSTGFWM1hkYTB0LzVQdGRHTVNFemN6ZW9aSFdrbkR0dys3OU9CaXR0UFBqDQpTaDFvRkR1UG8zNVI3ZVA2MjRsVUNjaC9JblpDcGhUYUx4OW9ETEdjYUszYWlsUTl3akJkS2RsQmw4S05LSVpwDQphMTNhUDVyblNtMkp2YSt0WHkveWkzQlNkczNkR0Q4SVRLWnlJLzZBRkh4R3ZPYnJESUJwbzRGRi96Y1dYVkRqDQpwYU9teHBsUnRNNEhpdG0rc1hHdmZxSmU0eDVEdU9YT25QclQzZEh2UlQ2dlNaVW9Lb2J4TXFtUlRPY3JPSVBhDQpFZU1wT29ic2hPUnVSbnRNRFl2dmdPM0Q2cDZpY2lEVzJWcDlONnJkTWRmT1dFUU44SlZXdkI3SXhSSGs5cUtKDQp2WU9XVmJjekF0MHFwTXZYRjNQWExqWmJVTTBrbk9kVUtJRWJxUDRZVWJnZHp4NlJ0Z2lpWTkzMEFqNnRBdGNlDQowZnBnTmx2ak1ScFNCdVdUbEFmTk5qRy9ZaG5kTXo5dUk2OFRNZkZwUjNQY2dWSXYzMGtydy85VnpvTGkyRHBlDQpvdzZEckdPNm9pK0RoTjc4UDRqWS9POVVjelpLMnJvWkwxT2k1UDBSSXhmMjNVWkM3eDFEbGNOM25CcjRzWVN2DQpyQng0Y0ZUTU5wd1UrbnpzSWk0ZGpjRkRLbUpkRU95ak1ua1AydjBMd2U3eXZLMDhwWmRFdSswemJycTE3a3VlDQpYcFhMYzdLNjhRQjE1eXh6R3lsVTVyUnd6bUMvWXNBVnlFNGVvR3U4UHhXeHJFUnZIYnk0QjhZUDB2QWZPcmFMDQpsS21YbEs0ZFRnPT0NCi0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0NCg==`

const validIntermediateAndCertRSA = `LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tDQpNSUlHTXpDQ0JCdWdBd0lCQWdJQkFUQU5CZ2txaGtpRzl3MEJBUXNGQURCZk1Rc3dDUVlEVlFRR0V3SlZVekVMDQpNQWtHQTFVRUNBd0NRMEV4RmpBVUJnTlZCQWNNRFZOaGJpQkdjbUZ1WTJselkyOHhEekFOQmdOVkJBb01Ca1J2DQpZMnRsY2pFYU1CZ0dBMVVFQXd3UlRtOTBZWEo1SUZSbGMzUnBibWNnUTBFd0hoY05NVFV3TnpFMk1EUXlOVEF6DQpXaGNOTWpVd056RXpNRFF5TlRBeldqQmZNUm93R0FZRFZRUUREQkZPYjNSaGNua2dWR1Z6ZEdsdVp5QkRRVEVMDQpNQWtHQTFVRUJoTUNWVk14RmpBVUJnTlZCQWNNRFZOaGJpQkdjbUZ1WTJselkyOHhEekFOQmdOVkJBb01Ca1J2DQpZMnRsY2pFTE1Ba0dBMVVFQ0F3Q1EwRXdnZ0lpTUEwR0NTcUdTSWIzRFFFQkFRVUFBNElDRHdBd2dnSUtBb0lDDQpBUUN3VlZENHBLN3o3cFhQcEpiYVoxSGc1ZVJYSWNhWXRiRlBDbk4waXF5OUhzVkVHbkVuNUJQTlNFc3VQK20wDQo1TjBxVlY3REdiMVNqaWxvTFhEMXFERHZoWFdrK2dpUzlwcHFQSFBMVlBCNGJ2enNxd0RZcnRwYnFrWXZPMFlLDQowU0wza3hQWFVGZGxrRmZndTB4amxjem0yUGhXRzNKZDhhQXRzcEwvTCtWZlBBMTNKVWFXeFNMcHVpMUluOHJoDQpnQXlRVEs2UTRPZjZHYkpZVG5BSGI1OVVvTFhTekI1QWZxaVVxNkw3bkVZWUtvUGZsUGJSQUlXTC9VQm0wYytIDQpvY21zNzA2UFlwbVBTMlJRdjNpT0dtbm45aEVWcDNQNmpxN1dBZXZiQTRhWUd4NUVzYlZ0WUFCcUpCYkZXQXV3DQp3VEdSWW16bjBNajBlVE1nZTl6dFlCMi8yc3hkVGU2dWhtRmdwVVhuZ0RxSkk1TzlOM3pQZnZsRUltQ2t5M0hNDQpqSm9MN2c1c21xWDlvMVArRVNMaDBWWnpoaDdJRFB6UVRYcGNQSVMvNnowbDIyUUdrSy8xTjFQYUFEYVVIZExMDQp2U2F2M3kyQmFFbVB2ZjJma1pqOHlQNWVZZ2k3Q3c1T05oSExEWUhGY2w5Wm0veXdtZHhISkVUejluZmdYbnNXDQpITnhEcXJrQ1ZPNDZyL3U2clNyVXQ2aHIzb2RkSkc4czhKbzA2ZWFydzZYVTNNek0rM2dpd2tLMFNTTTN1UlBxDQo0QXNjUjFUditFMzFBdU9BbWpxWVFvVDI5Yk1JeG9TemVsamovWW5lZHdqVzQ1cFd5YzNKb0hhaWJEd3ZXOVVvDQpHU1pCVnk0aHJNL0ZhN1hDV3YxV2ZITlcxZ0R3YUxZd0RubDVqRm1SQnZjZnVRSURBUUFCbzRINU1JSDJNSUdSDQpCZ05WSFNNRWdZa3dnWWFBRkhVTTFVM0U0V3lMMW52RmQrZFBZOGY0TzJoWm9XT2tZVEJmTVFzd0NRWURWUVFHDQpFd0pWVXpFTE1Ba0dBMVVFQ0F3Q1EwRXhGakFVQmdOVkJBY01EVk5oYmlCR2NtRnVZMmx6WTI4eER6QU5CZ05WDQpCQW9NQmtSdlkydGxjakVhTUJnR0ExVUVBd3dSVG05MFlYSjVJRlJsYzNScGJtY2dRMEdDQ1FEQ2VETGJlbUlUDQpTekFTQmdOVkhSTUJBZjhFQ0RBR0FRSC9BZ0VBTUIwR0ExVWRKUVFXTUJRR0NDc0dBUVVGQndNQ0JnZ3JCZ0VGDQpCUWNEQVRBT0JnTlZIUThCQWY4RUJBTUNBVVl3SFFZRFZSME9CQllFRkhlNDhoY0JjQXAwYlVWbFR4WGVSQTRvDQpFMTZwTUEwR0NTcUdTSWIzRFFFQkN3VUFBNElDQVFBV1V0QVBkVUZwd1JxK04xU3pHVWVqU2lrZU1HeVBac2NaDQpKQlVDbWhab0Z1ZmdYR2JMTzVPcGNSTGFWM1hkYTB0LzVQdGRHTVNFemN6ZW9aSFdrbkR0dys3OU9CaXR0UFBqDQpTaDFvRkR1UG8zNVI3ZVA2MjRsVUNjaC9JblpDcGhUYUx4OW9ETEdjYUszYWlsUTl3akJkS2RsQmw4S05LSVpwDQphMTNhUDVyblNtMkp2YSt0WHkveWkzQlNkczNkR0Q4SVRLWnlJLzZBRkh4R3ZPYnJESUJwbzRGRi96Y1dYVkRqDQpwYU9teHBsUnRNNEhpdG0rc1hHdmZxSmU0eDVEdU9YT25QclQzZEh2UlQ2dlNaVW9Lb2J4TXFtUlRPY3JPSVBhDQpFZU1wT29ic2hPUnVSbnRNRFl2dmdPM0Q2cDZpY2lEVzJWcDlONnJkTWRmT1dFUU44SlZXdkI3SXhSSGs5cUtKDQp2WU9XVmJjekF0MHFwTXZYRjNQWExqWmJVTTBrbk9kVUtJRWJxUDRZVWJnZHp4NlJ0Z2lpWTkzMEFqNnRBdGNlDQowZnBnTmx2ak1ScFNCdVdUbEFmTk5qRy9ZaG5kTXo5dUk2OFRNZkZwUjNQY2dWSXYzMGtydy85VnpvTGkyRHBlDQpvdzZEckdPNm9pK0RoTjc4UDRqWS9POVVjelpLMnJvWkwxT2k1UDBSSXhmMjNVWkM3eDFEbGNOM25CcjRzWVN2DQpyQng0Y0ZUTU5wd1UrbnpzSWk0ZGpjRkRLbUpkRU95ak1ua1AydjBMd2U3eXZLMDhwWmRFdSswemJycTE3a3VlDQpYcFhMYzdLNjhRQjE1eXh6R3lsVTVyUnd6bUMvWXNBVnlFNGVvR3U4UHhXeHJFUnZIYnk0QjhZUDB2QWZPcmFMDQpsS21YbEs0ZFRnPT0NCi0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0NCi0tLS0tQkVHSU4gQ0VSVElGSUNBVEUtLS0tLQ0KTUlJRlZ6Q0NBeitnQXdJQkFnSUJBekFOQmdrcWhraUc5dzBCQVFzRkFEQmZNUm93R0FZRFZRUUREQkZPYjNSaA0KY25rZ1ZHVnpkR2x1WnlCRFFURUxNQWtHQTFVRUJoTUNWVk14RmpBVUJnTlZCQWNNRFZOaGJpQkdjbUZ1WTJseg0KWTI4eER6QU5CZ05WQkFvTUJrUnZZMnRsY2pFTE1Ba0dBMVVFQ0F3Q1EwRXdIaGNOTVRVd056RTJNRFF5TlRVdw0KV2hjTk1UWXdOekUxTURReU5UVXdXakJnTVJzd0dRWURWUVFEREJKelpXTjFjbVV1WlhoaGJYQnNaUzVqYjIweA0KQ3pBSkJnTlZCQVlUQWxWVE1SWXdGQVlEVlFRSERBMVRZVzRnUm5KaGJtTnBjMk52TVE4d0RRWURWUVFLREFaRQ0KYjJOclpYSXhDekFKQmdOVkJBZ01Ba05CTUlJQklqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FROEFNSUlCQ2dLQw0KQVFFQW1MWWlZQ1RBV0pCV0F1eFpMcVZtVjRGaVVkR2dFcW9RdkNiTjczekYvbVFmaHEwQ0lUbzZ4U3hzMVFpRw0KRE96VXRrcHpYenppU2o0SjUrZXQ0SmtGbGVlRUthTWNIYWRlSXNTbEhHdlZ0WER2OTNvUjN5ZG1mWk8rVUxSVQ0KOHhIbG9xY0xyMUtyT1AxZGFMZmRNUmJhY3RkNzVVUWd2dzlYVHNkZU1WWDVBbGljU0VOVktWK0FRWHZWcHY4UA0KVDEwTVN2bEJGYW00cmVYdVkvU2tlTWJJYVc1cEZ1NkFRdjNabWZ0dDJ0YTBDQjlrYjFtWWQrT0tydThIbm5xNQ0KYUp3NlIzR2hQMFRCZDI1UDFQa2lTeE0yS0dZWlprMFcvTlpxTEs5L0xURktUTkN2N1ZqQ2J5c1ZvN0h4Q1kwYg0KUWUvYkRQODJ2N1NuTHRiM2Fab2dmdmE0SFFJREFRQUJvNElCR3pDQ0FSY3dnWWdHQTFVZEl3U0JnREIrZ0JSMw0KdVBJWEFYQUtkRzFGWlU4VjNrUU9LQk5lcWFGanBHRXdYekVMTUFrR0ExVUVCaE1DVlZNeEN6QUpCZ05WQkFnTQ0KQWtOQk1SWXdGQVlEVlFRSERBMVRZVzRnUm5KaGJtTnBjMk52TVE4d0RRWURWUVFLREFaRWIyTnJaWEl4R2pBWQ0KQmdOVkJBTU1FVTV2ZEdGeWVTQlVaWE4wYVc1bklFTkJnZ0VCTUF3R0ExVWRFd0VCL3dRQ01BQXdIUVlEVlIwbA0KQkJZd0ZBWUlLd1lCQlFVSEF3SUdDQ3NHQVFVRkJ3TUJNQTRHQTFVZER3RUIvd1FFQXdJRm9EQXVCZ05WSFJFRQ0KSnpBbGdoSnpaV04xY21VdVpYaGhiWEJzWlM1amIyMkNDV3h2WTJGc2FHOXpkSWNFZndBQUFUQWRCZ05WSFE0RQ0KRmdRVURQRDRDYVhSYnU1UUJiNWU4eThvZHZUcVc0SXdEUVlKS29aSWh2Y05BUUVMQlFBRGdnSUJBSk95bG1jNA0KbjdKNjRHS3NQL3hoVWRLS1Y5L0tEK3VmenBLYnJMSW9qV243clR5ZTcwdlkwT2pRRnVPWGM1NHlqTVNJTCsvNQ0KbWxOUTdZL2ZKUzh4ZEg3OUVSKzRuV011RDJlY2lMbnNMZ2JZVWs0aGl5Ynk4LzVWKy9ZcVBlQ3BQQ242VEpSSw0KYTBFNmxWL1VqWEpkcmlnSnZKb05PUjhaZ3RFWi9RUGdqSkVWVXNnNDdkdHF6c0RwZ2VTOGRjanVNV3BaeFAwMg0KcWF2RkxEalNGelZIKzJENk90eTFEUXBsbS8vM1hhUlhoMjNkT0NQOHdqL2J4dm5WVG9GV3Mrek80dVQxTEYvUw0KS1hDTlFvZWlHeFdIeXpyWEZWVnRWbkM5RlNOejBHZzIvRW0xdGZSZ3ZoVW40S0xKY3ZaVzlvMVI3VlZDWDBMMQ0KMHgwZnlLM1ZXZVdjODZhNWE2ODFhbUtaU0ViakFtSVZaRjl6T1gwUE9EQzhveSt6cU9QV2EwV0NsNEs2ekRDNg0KMklJRkJCTnk1MFpTMmlPTjZSWTZtRTdObUE3OGdja2Y0MTVjcUlWcmxvWUpiYlREZXBmaFRWMjE4U0xlcHBoNA0KdUdiMi9zeGtsZkhPWUUrcnBIY2lpYld3WHJ3bE9ESmFYdXpYRmhwbFVkL292ZHVqQk5BSUhrQmZ6eStZNnoycw0KYndaY2ZxRDROSWIvQUdoSXlXMnZidnU0enNsRHAxTUVzTG9hTytTemlyTXpreU1CbEtSdDEyMHR3czRFa1VsbQ0KL1FoalNVb1pwQ0FzeTVDL3BWNCtieDBTeXNOZC9TK2tLYVJaYy9VNlkzWllCRmhzekxoN0phTFhLbWs3d0huRQ0KcmdnbTZvejRML0d5UFdjL0ZqZm5zZWZXS00yeUMzUURoanZqDQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tDQo=`

const signedRSARootTemplate = `{"signed":{"_type":"Root","consistent_snapshot":false,"expires":"2016-07-16T23:34:13.389129622-07:00","keys":{"1fc4fdc38f66558658c5c59b67f1716bdc6a74ef138b023ae5931db69f51d670":{"keytype":"ecdsa","keyval":{"private":null,"public":"MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE8nIgzLigo5D47dWQe1IUjzHXxvyx0j/OL16VQymuloWsgVDxxT6+mH3CeviMAs+/McnEPE9exnm6SQGR5x3XMw=="}},"23c29cc372109c819e081bc953b7657d05e3f968f03c21d0d75ea457590f3d14":{"keytype":"ecdsa","keyval":{"private":null,"public":"MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEClUFVWkc85OQScfTQRS02VaLIEaeCmxdwYS/hcTLVoTxlFfRfs7HyalTwXGAGO79XZZS+koE6s8D0xGcCJQkLQ=="}},"49cf5c6404a35fa41d5a5aa2ce539dfee0d7a2176d0da488914a38603b1f4292":{"keytype":"rsa-x509","keyval":{"private":null,"public":"{{.RootPem}}"}},"e3a5a4fdaf11ea1ec58f5efed6f3639b39cd4cfa1418c8b55c9a8c2447ace5d9":{"keytype":"ecdsa","keyval":{"private":null,"public":"MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEgl3rzMPMEKhS1k/AX16MM4PdidpjJr+z4pj0Td+30QnpbOIARgpyR1PiFztU8BZlqG3cUazvFclr2q/xHvfrqw=="}}},"roles":{"root":{"keyids":["49cf5c6404a35fa41d5a5aa2ce539dfee0d7a2176d0da488914a38603b1f4292"],"threshold":1},"snapshot":{"keyids":["23c29cc372109c819e081bc953b7657d05e3f968f03c21d0d75ea457590f3d14"],"threshold":1},"targets":{"keyids":["1fc4fdc38f66558658c5c59b67f1716bdc6a74ef138b023ae5931db69f51d670"],"threshold":1},"timestamp":{"keyids":["e3a5a4fdaf11ea1ec58f5efed6f3639b39cd4cfa1418c8b55c9a8c2447ace5d9"],"threshold":1}},"version":2},"signatures":[{"keyid":"49cf5c6404a35fa41d5a5aa2ce539dfee0d7a2176d0da488914a38603b1f4292","method":"rsapss","sig":"YlZwtCj028Xc23+KHfj6govFEY6hMbBXO5HT20F0I5ZeIPb1l7OmkjEiwp9ZHusClY+QeqiP1CFh\n/AfCbv4tLanqMkXPtm8UJJ1hMZVq86coieB32PQDj9k6x1hErHzvPUbOzTRW2BQkFFMZFkLDAd06\npH8lmxyPLOhdkVE8qIT7sBCy/4bQIGfvEX6yCDz84MZdcLNX5B9mzGi9A7gDloh9IEZxA8UgoI18\nSYpv/fYeSZSqM/ws2G+kiELGgTWhcZ+gOlF7ArM/DOlcC/NYqcvY1ugE6Gn7G8opre6NOofdRp3w\n603A2rMMvYTwqKLY6oX/d+07A2+WGHXPUy5otCAybWOw2hIZ35Jjmh12g6Dc6Qk4K2zXwAgvWwBU\nWlT8MlP1Tf7f80jnGjh0aARlHI4LCxlYU5L/pCaYuHgynujvLuzoOuiiPfJv7sYvKoQ8UieE1w//\nHc8E6tWtV5G2FguKLurMoKZ9FBWcanDO0fg5AWuG3qcgUJdvh9acQ33EKer1fqBxs6LSAUWo8rDt\nQkg+b55AW0YBukAW9IAfMySQGAS2e3mHZ8nK/ijaygCRu7/P+NgKY9/zpmfL2xgcNslLcANcSOOt\nhiJS6yqYM9i9G0af0yw/TxAT4ntwjVm8u52UyR/hXIiUc/mjZcYRbSmJOHws902+i+Z/qv72knk="}]}`

// TestValidateRoot through the process of initializing a repository and makes
// sure the repository looks correct on disk.
// We test this with both an RSA and ECDSA root key
func TestValidateRoot(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	validateRootSuccessfully(t, data.RSAKey)
}

func validateRootSuccessfully(t *testing.T, rootType data.KeyAlgorithm) {
	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	defer os.RemoveAll(tempBaseDir)

	assert.NoError(t, err, "failed to create a temporary directory: %s", err)

	gun := "docker.com/notary"

	ts, mux := createTestServer(t)
	defer ts.Close()

	repo, err := NewNotaryRepository(tempBaseDir, gun, ts.URL, http.DefaultTransport, passphraseRetriever)
	assert.NoError(t, err, "error creating repository: %s", err)

	rootKeyID, err := repo.KeyStoreManager.GenRootKey(rootType.String())
	assert.NoError(t, err, "error generating root key: %s", err)

	rootCryptoService, err := repo.KeyStoreManager.GetRootCryptoService(rootKeyID)
	assert.NoError(t, err, "error retrieving root key: %s", err)

	err = repo.Initialize(rootCryptoService)
	assert.NoError(t, err, "error creating repository: %s", err)

	// Initialize is supposed to have created new certificate for this repository
	// Lets check for it and store it for later use
	allCerts := repo.KeyStoreManager.TrustedCertificateStore().GetCertificates()
	assert.Len(t, allCerts, 1)

	// Now test ListTargets. In preparation, we need to expose some signed
	// metadata files on the internal HTTP server.
	var tempKey data.TUFKey
	json.Unmarshal([]byte(timestampECDSAKeyJSON), &tempKey)

	repo.KeyStoreManager.NonRootKeyStore().AddKey(filepath.Join(filepath.FromSlash(gun), tempKey.ID()), "root", &tempKey)

	// Because ListTargets will clear this
	savedTUFRepo := repo.tufRepo

	rootJSONFile := filepath.Join(tempBaseDir, "tuf", filepath.FromSlash(gun), "metadata", "root.json")
	rootFileBytes, err := ioutil.ReadFile(rootJSONFile)

	signedTargets, err := savedTUFRepo.SignTargets("targets", data.DefaultExpires("targets"), nil)
	assert.NoError(t, err)

	signedSnapshot, err := savedTUFRepo.SignSnapshot(data.DefaultExpires("snapshot"), nil)
	assert.NoError(t, err)

	signedTimestamp, err := savedTUFRepo.SignTimestamp(data.DefaultExpires("timestamp"), nil)
	assert.NoError(t, err)

	mux.HandleFunc("/v2/docker.com/notary/_trust/tuf/root.json", func(w http.ResponseWriter, r *http.Request) {
		assert.NoError(t, err)
		fmt.Fprint(w, string(rootFileBytes))
	})

	mux.HandleFunc("/v2/docker.com/notary/_trust/tuf/timestamp.json", func(w http.ResponseWriter, r *http.Request) {
		timestampJSON, _ := json.Marshal(signedTimestamp)
		fmt.Fprint(w, string(timestampJSON))
	})

	mux.HandleFunc("/v2/docker.com/notary/_trust/tuf/snapshot.json", func(w http.ResponseWriter, r *http.Request) {
		snapshotJSON, _ := json.Marshal(signedSnapshot)
		fmt.Fprint(w, string(snapshotJSON))
	})

	mux.HandleFunc("/v2/docker.com/notary/_trust/tuf/targets.json", func(w http.ResponseWriter, r *http.Request) {
		targetsJSON, _ := json.Marshal(signedTargets)
		fmt.Fprint(w, string(targetsJSON))
	})

	_, err = repo.ListTargets()
	assert.NoError(t, err)

	//
	// Test TOFUS logic. We remove all certs and expect a new one to be added after ListTargets
	//
	err = repo.KeyStoreManager.TrustedCertificateStore().RemoveAll()
	assert.NoError(t, err)
	assert.Len(t, repo.KeyStoreManager.TrustedCertificateStore().GetCertificates(), 0)

	// This list targets is expected to succeed and the certificate store to have the new certificate
	_, err = repo.ListTargets()
	assert.NoError(t, err)
	assert.Len(t, repo.KeyStoreManager.TrustedCertificateStore().GetCertificates(), 1)

	//
	// Test certificate mismatch logic. We remove all certs, add a different cert to the
	// same CN, and expect ValidateRoot to fail
	//

	// First, remove all certs
	err = repo.KeyStoreManager.TrustedCertificateStore().RemoveAll()
	assert.NoError(t, err)
	assert.Len(t, repo.KeyStoreManager.TrustedCertificateStore().GetCertificates(), 0)

	// Add a previously generated certificate with CN=docker.com/notary
	err = repo.KeyStoreManager.TrustedCertificateStore().AddCertFromFile("../fixtures/self-signed_docker.com-notary.crt")
	assert.NoError(t, err)

	// This list targets is expected to fail, since there already exists a certificate
	// in the store for the dnsName docker.com/notary, so TOFUS doesn't apply
	_, err = repo.ListTargets()
	if assert.Error(t, err, "An error was expected") {
		assert.Equal(t, err, &keystoremanager.ErrValidationFail{Reason: "failed to validate integrity of roots"})
	}

	//
	// Test CA fallback logic.
	// If we already have a certificate that matches a CN and get presented a new
	// certificate for that same CN and that certificate is signed by a CA we
	// know about, then we trust it but don't remove the previous pinning.
	//

	// First, remove all certs and trusted CAs
	err = repo.KeyStoreManager.TrustedCertificateStore().RemoveAll()
	assert.NoError(t, err)
	assert.Len(t, repo.KeyStoreManager.TrustedCertificateStore().GetCertificates(), 0)
	err = repo.KeyStoreManager.TrustedCAStore().RemoveAll()
	assert.NoError(t, err)
	assert.Len(t, repo.KeyStoreManager.TrustedCAStore().GetCertificates(), 0)

	// Add a trusted root-ca
	err = repo.KeyStoreManager.TrustedCAStore().AddCertFromFile("../fixtures/root-ca.crt")
	assert.NoError(t, err)

	// Add a previously generated certificate with CN=secure.example.com
	err = repo.KeyStoreManager.TrustedCertificateStore().AddCertFromFile("../fixtures/self-signed_secure.example.com.crt")
	assert.NoError(t, err)

	// Create a manual rootSigner with a valid intermediate + leaf cert
	var testSignedRoot data.Signed
	var signedRootBytes bytes.Buffer

	// Execute our template
	templ, _ := template.New("SignedRSARootTemplate").Parse(signedRSARootTemplate)
	templ.Execute(&signedRootBytes, SignedRSARootTemplate{RootPem: validIntermediateAndCertRSA})

	// Unmarshal our signedroot
	json.Unmarshal(signedRootBytes.Bytes(), &testSignedRoot)

	//
	// This call to ValidateRoot will suceed in getting to the TUF validation, since
	// we are usign a valid PEM encoded certificate chain of intermediate + leaf cert
	// that are signed by a trusted root authority and the leaf cert has a correct CN.
	// It will, however, fail to validate, because this is has an invalid TUF signature
	//
	err = repo.KeyStoreManager.ValidateRoot(&testSignedRoot, "secure.example.com")
	if assert.Error(t, err, "An error was expected") {
		assert.Equal(t, err, &keystoremanager.ErrValidationFail{Reason: "failed to validate integrity of roots"})
	}
}

func TestValidateRootWithInvalidData(t *testing.T) {
	var testSignedRoot data.Signed
	var signedRootBytes bytes.Buffer

	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	defer os.RemoveAll(tempBaseDir)
	assert.NoError(t, err, "failed to create a temporary directory: %s", err)

	// Create a FileStoreManager
	keyStoreManager, err := keystoremanager.NewKeyStoreManager(tempBaseDir, passphraseRetriever)
	assert.NoError(t, err)

	// Execute our template
	templ, _ := template.New("SignedRSARootTemplate").Parse(signedRSARootTemplate)
	templ.Execute(&signedRootBytes, SignedRSARootTemplate{RootPem: validPEMEncodedRSARoot})

	// Unmarshal our signedroot
	json.Unmarshal(signedRootBytes.Bytes(), &testSignedRoot)

	//
	// This call to ValidateRoot will succeed since we are usign a valid PEM
	// encoded certificate, and have no other certificates for this CN
	//
	err = keyStoreManager.ValidateRoot(&testSignedRoot, "docker.com/notary")
	assert.NoError(t, err)

	//
	// This call to ValidateRoot will fail since we are passing in a dnsName that
	// doesn't match the CN of the certificate.
	//
	err = keyStoreManager.ValidateRoot(&testSignedRoot, "diogomonica.com/notary")
	if assert.Error(t, err, "An error was expected") {
		assert.Equal(t, err, &keystoremanager.ErrValidationFail{Reason: "unable to retrieve valid leaf certificates"})
	}

	//
	// This call to ValidateRoot will fail since we are passing an unparsable RootSigned
	//
	// Execute our template deleting the old buffer first
	signedRootBytes.Reset()
	templ, _ = template.New("SignedRSARootTemplate").Parse(signedRSARootTemplate)
	templ.Execute(&signedRootBytes, SignedRSARootTemplate{RootPem: "------ ABSOLUTELY NOT A PEM -------"})
	// Unmarshal our signedroot
	json.Unmarshal(signedRootBytes.Bytes(), &testSignedRoot)

	err = keyStoreManager.ValidateRoot(&testSignedRoot, "docker.com/notary")
	assert.Error(t, err, "illegal base64 data at input byte")

	//
	// This call to ValidateRoot will fail since we are passing an invalid PEM cert
	//
	// Execute our template deleting the old buffer first
	signedRootBytes.Reset()
	templ, _ = template.New("SignedRSARootTemplate").Parse(signedRSARootTemplate)
	templ.Execute(&signedRootBytes, SignedRSARootTemplate{RootPem: "LS0tLS1CRUdJTiBDRVJU"})
	// Unmarshal our signedroot
	json.Unmarshal(signedRootBytes.Bytes(), &testSignedRoot)

	err = keyStoreManager.ValidateRoot(&testSignedRoot, "docker.com/notary")
	if assert.Error(t, err, "An error was expected") {
		assert.Equal(t, err, &keystoremanager.ErrValidationFail{Reason: "unable to retrieve valid leaf certificates"})
	}

	//
	// This call to ValidateRoot will fail since we are passing only CA certificate
	// This will fail due to the lack of a leaf certificate
	//
	// Execute our template deleting the old buffer first
	signedRootBytes.Reset()
	templ, _ = template.New("SignedRSARootTemplate").Parse(signedRSARootTemplate)
	templ.Execute(&signedRootBytes, SignedRSARootTemplate{RootPem: validCAPEMEncodeRSARoot})
	// Unmarshal our signedroot
	json.Unmarshal(signedRootBytes.Bytes(), &testSignedRoot)

	err = keyStoreManager.ValidateRoot(&testSignedRoot, "docker.com/notary")
	if assert.Error(t, err, "An error was expected") {
		assert.Equal(t, err, &keystoremanager.ErrValidationFail{Reason: "unable to retrieve valid leaf certificates"})
	}
}
