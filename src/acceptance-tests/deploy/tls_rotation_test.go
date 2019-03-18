package deploy_test

import (
	"fmt"
	"time"

	"github.com/cloudfoundry-incubator/consul-release/src/acceptance-tests/testing/consulclient"
	"github.com/cloudfoundry-incubator/consul-release/src/acceptance-tests/testing/helpers"
	"github.com/pivotal-cf-experimental/bosh-test/bosh"
	"github.com/pivotal-cf-experimental/destiny/ops"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	newCACert = `-----BEGIN CERTIFICATE-----
MIIE5jCCAs6gAwIBAgIBATANBgkqhkiG9w0BAQsFADATMREwDwYDVQQDEwhjb25z
dWxDQTAeFw0xOTAzMTgwODQyNDBaFw0yMDA5MTgwODQyNDBaMBMxETAPBgNVBAMT
CGNvbnN1bENBMIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAogmIAzeO
q59mleCRHTRMSTV6JuxiL9XR+njqEGdl7mzRcWOg3E70E+3mrwCtHiISkREN8Ymc
H/YeWoPOfGgEfvIeXhFplxlB5V1LwXyDaDovuA72e1XYkubAwV5NR9C/g4f/iGO7
1NX3tjjBHvRET3MtY1E8O5faC8w1I7+bo0Dy76tNCRjgGB50T9n5jcXmN+uQgmyW
cZET9lh9PKMVt5ecxFr+Kh+YzXckaDzVAGoYvGrEZ+jlevOe25c1Qv4eKB6k+u97
Z0H+cT9hBmO8r2oOuixc2RS5Jx5yY8MSMOPz859vZChoSoPEu20upt19raiJxdgM
06euXJxuuedzX0cg0cZL0ljqHQ6RO8LCaxpJTjVWRhWu6c+s0NYL4eO4i3alQUet
IAOhWJvs1qLF9r7vUsofSZWKTx8U8VrED69bzvQOfl8qrWTEk87BfDjuHn2O0RZc
r3s9fXgQ8n0VqZZkatYev4LPF8TW9VTOTbE4dIIElaRKuFgFddFG0VXD9Tvzd0IN
tot95KTgo1ojMgvpnPLLLn2toVxi9mw15BXQqiaBEExeKUjkKSE4b5F/1yxtCs1J
QdATXG4HgRoZTQVnQtPHRu3SP7xP9cefYbyLK3gIyxV9J0om92jGB2cxUOr2hSEb
31OrvuaOlVdN9w46kMtpDb6IWgDGxhLuQ40CAwEAAaNFMEMwDgYDVR0PAQH/BAQD
AgEGMBIGA1UdEwEB/wQIMAYBAf8CAQAwHQYDVR0OBBYEFLT41TQK6lv+NTee7DSi
P/LLXroSMA0GCSqGSIb3DQEBCwUAA4ICAQAkyFrq40pXy5Rkc9VniKiSXJ3LhWZx
DD+6HMdBeT8joIpM6GvSmge1L/RBuGGPBLw62z2ijFOjwognMsX5fJQ+lRfMQGb5
8edAfYlikGwK5+8v/FY7u6G1WQMPMbtsqx8mtMAfqliZhkwJIhpwemoaObVYMIhy
o7S2nNiVgFWSCAQCtPyv94v3sRz3rSej9TeD3cK3u7lA2QufuLr+WQoSR/4Zrt6o
qH9dgFC/0sxCAPDc5eY4WuVMKpu8C+iEk41X8r0A4fGy7DrVhmLYhaFx3Vs5QpuR
CSapGEi6gicwYk7+28DEWiNRAxMYgadDoKEP8rXjzhcUrhhuNWVHa86JZVyPtaEo
RLece56yFl+QOOKnHCBLO3mh8CkYudZ3A6Mhh6tyARaBipK34eA4dD/WD9DreZVN
5gwXI+h1bMeHalj3/vfP19K0ERlLN5zcWZPMrCkWoH+W0qEMujKha3NRQcOW1tjd
y59j/0x1kFSs1KsuLc11MhZa7hEKhUcK7jz1KDvz4wLsmcUmW+O1Zsxh0vLNKOkS
0Ob6GxEnstKwBWGuzDdAsfrtCLehhggWwEQlBPGtOIU/bq+LEnmFvXriHd8KE7wZ
o0sPOI2pMH+mfIrOGBKQ7D7PY9KFvBmkIT4UNnGThshMwoOfStg+4YyqJl6LIagW
KRsBJmUp9pG+lQ==
-----END CERTIFICATE-----`
	newAgentCert = `-----BEGIN CERTIFICATE-----
MIIEOTCCAiGgAwIBAgIRAOPA0aT1Yjo4aP7gKV9HuvkwDQYJKoZIhvcNAQELBQAw
EzERMA8GA1UEAxMIY29uc3VsQ0EwHhcNMTkwMzE4MDg0MjQxWhcNMjAwOTE4MDg0
MjQwWjAXMRUwEwYDVQQDEwxjb25zdWwgYWdlbnQwggEiMA0GCSqGSIb3DQEBAQUA
A4IBDwAwggEKAoIBAQCgEOSpqWGvWJG/xQthF+3w3u6szKNzGrXO0nUP1+/5j/nn
LTCYhM3noCdFq2KPut+N5z4gljcsGfkesgQp/RyrPznb3GGuD6FGclWuejAjdjnB
H830u7JIdDGMDVQwrDTJKkzej9fsevcOFSm+SHtuQMS7KV+Q87qOwEnWaY2bRGFu
qZiOX0hPl8eX42WWWp3/IyhGsJ5tJ+lkYjSN9wGYG3bfZuKhWcPazDFXJjGIjkat
cy6olbtAUMl1ij4/Ozr8IlRXplSbMlIdTrKMmjThHR3HjalRZczTjZNjmNvuow7+
SrcFx0iJJ6hEXhJCVsn6J0sABY/xzvmhgsvtolVXAgMBAAGjgYMwgYAwDgYDVR0P
AQH/BAQDAgO4MB0GA1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjAdBgNVHQ4E
FgQUsYQHWX6cHvxLF06f4qz3DYu6inowHwYDVR0jBBgwFoAUtPjVNArqW/41N57s
NKI/8steuhIwDwYDVR0RBAgwBocEfwAAATANBgkqhkiG9w0BAQsFAAOCAgEAIbWP
dxZBhyMm5ZN096xJHdpjmcHE/+E6Kh0AgUzFEOisMOpONPVCqJLi7L/5ygXHvAKn
tO/+Co+Xq4mUCFrqDZVtJNPnuYaZDB6rVcqg/Gh5I/aH+28rkTHGvvR+80gyRSKh
CB0T/w+kG19m+WV/gy4kIkaumMz8gehERHlVsYthih40Oc0Je2/hxHoCzQ6mr4Vl
W2ALma5nfdBG70KDEuzFXPfxPcZnMAqIVV4orbzpTyLIrJORBpUwqzeNCjim7GWZ
wEDpGyf25GBxY99vRQbJUUgCVKwulOXjlILyys9h/s86PgQE7cYs4HphKJN4nay2
PMFVyUA6qxbfA0LY+MFZqW5M/172+O7+0Ov9W2ifJK3dlKCVRCxDSOcacVPHWGy6
0nFzZBAemtA3a2MdtW9au6x1gB+ml0D/41T/TeJfZlnag3FCdPc1Ow8qEVOLn5rk
dg5FWmVy0hw7LYU7pFFdR36AgezUkILffwwqvc43QLqRJwLgEwKEmL+vGntMy3rc
cOEolUfGtUelZcgXtDSouc/Lp52kLlutPaLnDK7oaElrelxhw+bemFhPqdH4JHiB
MlhghwqV2gTgTzpLDb+H0fozCiuYDqkQLh32E+JIFJE4wjvhTFrFTs8cqX0T0olM
a+S6vfZ0xWZJTuQ1c0GO9qLjsuhLtgpE4wAdi0A=
-----END CERTIFICATE-----`
	newAgentKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAoBDkqalhr1iRv8ULYRft8N7urMyjcxq1ztJ1D9fv+Y/55y0w
mITN56AnRatij7rfjec+IJY3LBn5HrIEKf0cqz8529xhrg+hRnJVrnowI3Y5wR/N
9LuySHQxjA1UMKw0ySpM3o/X7Hr3DhUpvkh7bkDEuylfkPO6jsBJ1mmNm0RhbqmY
jl9IT5fHl+Nlllqd/yMoRrCebSfpZGI0jfcBmBt232bioVnD2swxVyYxiI5GrXMu
qJW7QFDJdYo+Pzs6/CJUV6ZUmzJSHU6yjJo04R0dx42pUWXM042TY5jb7qMO/kq3
BcdIiSeoRF4SQlbJ+idLAAWP8c75oYLL7aJVVwIDAQABAoIBAADUlZsbudoDB56L
EygJy744KdzTovVx6geMp/bRE/mjeZRtc5cW+Up+VjXSUcyVF5vQ202n+dlMuTIj
kkn1ejSZO1+coRUaF6gZ57/j+mP0tQ//bj4ayy39DFTBrPIjspJomcV90Yz0hluY
WIMYNSu5QkwGRuyllw4uiTOHkKiBjDXcmkpkKZUOvEFPMPqPJb2UkNou9yCs4C92
6fpspfe29LG67FZ/pjAQORcjnN/iiCHAkeIEJYh7Zk59xoRd+VnCKjbJlrsZOtXq
bga9VNlDFfMxUXvmIDc6/yJK3frfFxnr9VU7D6SkVUisJG0sTzfRnek0TR8BEyU7
W3SaGUECgYEAw9W1rzQpaRlrIoME6ENcdC8mp8yR3v9saIomzBOzLMHwVSUoOAz9
hdtJuJbiorZMq5uYlGCsw8saTtwAQs2lfbZCyO5Q2GWJ5iY4MogOFXkT5vNMZ344
xCaKwFeIQHYMPmdTvL5qLeZnTEfZKvd0du4iUn5wDx7IYzPvACLPIQ8CgYEA0T38
1QcLuibgXMP332T36gR2JxDYfXC93fYEHCz+Q/+jS2AuDOWwDga605WR2n5ftfY4
PfsWe4COu29P/9qgfjow4gkbHlPOHk157pKVd858QSkioUjts1NASowQZvZHBgoe
sRnj+brdBuZc7FNdgyN8I476Jv6Mf5vhlZx6dzkCgYBdgJ1VwOcAulUvzjS6nOb4
xkaDmaYQPg5Jv6SUjddfyF1ymeIhGPq3PaNuUgR1werLiOgJ+Dqk5UVzX5F0U/Hv
GuW6QCczmw+DZr4wSkvHLt80xve09kwuQ2S+P0zb0kE4Tmdp19Skg7zQbAGhhTMD
UeHrV1kzrvPogbRccUJKOwKBgBroidcrbMqnrTrAyOOlrGwf3sHvXKflE8WzmZu1
/YzpFyreV425DAcBvozvMy6SCeTwoRL3c1C2m6RnEDaq+vDAswCegypHRL6I4CFa
IHajyz7l91oectMY5a+wi3tyOHgCXSgRWEwJR9tXTKPnpKL0sUYxYOIa4h6XAU+o
K+ehAoGBAMNxdxnohkVOCIyN+Nj+CI6V56lH8I8LAYAnnyYTdLrjHcqWJLpb527f
NEXvO0Alo+nbvrtq1VLYRxZD67GkLZC23yTu3K+RY1cjr92sbV/vYrBfONW3SLeA
0HL0aeVKo6Vi67aLJbDxzjb6cZIssMHe2ptiMMiAVz2A6LNehZpU
-----END RSA PRIVATE KEY-----`
	newServerCert = `-----BEGIN CERTIFICATE-----
MIIELzCCAhegAwIBAgIQBoPpQrQksjOzIVToINOpojANBgkqhkiG9w0BAQsFADAT
MREwDwYDVQQDEwhjb25zdWxDQTAeFw0xOTAzMTgwODQyNDFaFw0yMDA5MTgwODQy
NDBaMCExHzAdBgNVBAMTFnNlcnZlci5kYzEuY2YuaW50ZXJuYWwwggEiMA0GCSqG
SIb3DQEBAQUAA4IBDwAwggEKAoIBAQC2JlcZFwPXZtFwVB1ChvgLyG4iBfFggjdE
L4Qf65FRSqfChObkoqTcDOMrGts0hid87x8BCTcG8qStLAbJPRiYzcK6aBgSSrKP
OU/sVpm8LylY0fjx9weTCTtixBUTcHfGRZL2dEFGYsaWXKzfWd4GU8WKXuZKHKB9
w7ZW/BgRNrxn7NHGRNiH+SGuDk/sMXorsAypF8XLn7EytLnms+AKEHqeNSIDdg58
Dop0nl/Efd2rO0RpicxEwJesV0c70kmiEpMUCP6kiCe6wc3JuwdFSevf3wsKFcrP
dEsqSszSZ3FmU0/9qeu23VhIkmbQBatg/gGgW7Tpk+UhrlUgt2o3AgMBAAGjcTBv
MA4GA1UdDwEB/wQEAwIDuDAdBgNVHSUEFjAUBggrBgEFBQcDAQYIKwYBBQUHAwIw
HQYDVR0OBBYEFGOz5aBT/QYnDuwmTyd6FRy1RBofMB8GA1UdIwQYMBaAFLT41TQK
6lv+NTee7DSiP/LLXroSMA0GCSqGSIb3DQEBCwUAA4ICAQAkTOUzHnJvapaKZIf0
aDCjC/LArl5WOlfxKPZaUcnbunjMAtgoXVBd9AY4wtIfo8HxX38EUV/nQ2NxFupD
waSR0qMaM1gqSrlNfik/IITu4pdLbmozfRQOaT9a6PJWceBiwx2bLXj604yDzrPr
zBdIGWt15Yu60eCyIpgZCf+dE3tX+z1zaYHHAu6Ca79EWO3raJDm4PHE2PNN5JSU
5Vbm1zqSZKn3U65rRtqq4uF00trohhdP6+N048Sqrj/lC/gXsQJWaD+P9HH68FR/
jNmUANHBoGeVvEtl1+kAhA7d8e37ZHpmsfGIcBCMeuZejkXQVa96uNVKkib3NV2X
3JF6qktkv2UxZJ79eBBMlomZZkMwopYxa/mlCL26hhQDgjVJ/fygwH1qfZG+SjEK
1xbg7Lj8r3uaB+WTcm6Zzw0t7cm0i26omJI9cRD6FFZDGaCANkeEMxQqsccyEiKL
jNYVv1ZSP6IEFtzLL8bLz194jZ6fbkBe+pauyS2SihHHzcWknLCGK8S1joZsIvx2
uHrsT4rBte1C6Gb8WEeTMIcivPR0n/JpBu3Nqvs6s3ct8HdkaOaPJQfRZaxKUVhH
W9xCh2w0GEEygQNse5QPpirJZEPSRyw8jzuAVP4u1m9q7J9hIj6IWoW7IjkcnqFI
kP3EWl/e7dqf2O2+GHGAjuH6nQ==
-----END CERTIFICATE-----`
	newServerKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAtiZXGRcD12bRcFQdQob4C8huIgXxYII3RC+EH+uRUUqnwoTm
5KKk3AzjKxrbNIYnfO8fAQk3BvKkrSwGyT0YmM3CumgYEkqyjzlP7FaZvC8pWNH4
8fcHkwk7YsQVE3B3xkWS9nRBRmLGllys31neBlPFil7mShygfcO2VvwYETa8Z+zR
xkTYh/khrg5P7DF6K7AMqRfFy5+xMrS55rPgChB6njUiA3YOfA6KdJ5fxH3dqztE
aYnMRMCXrFdHO9JJohKTFAj+pIgnusHNybsHRUnr398LChXKz3RLKkrM0mdxZlNP
/anrtt1YSJJm0AWrYP4BoFu06ZPlIa5VILdqNwIDAQABAoIBABEGuHGB8zv5Qm6L
jkifsSP40kKf55Yr1dqrzl/ldynwHopSPfr1MQ/YrItk8USRnbNR7sX8BIbDu5zs
Vp4M9fWilicyf72F+fblVpEy7x/mEKlaxzhm7PBTlpd+2LF+e9OuwTQEXe2kYgfA
FfCyx8wznG14vXIEBwR6fNrYqa9CFC9z7fW3PzBo3B/QLfzoJgUN31N2uKPAyMY2
0WI9Pk/DkpmHASGOxL0t681RZWHLFQ8cq0tilYJvoi10Pk+D2eXYu9rDEDdACZI4
Mt4GCjQppvjoxYGZVpqu9CHOWowWbd5SUQeKuDO/xaGKJTgoMdSHOdHLMRSPBbxn
DdrJn4ECgYEA2Tnwsoc6z0EnimPkIEKs3Y0/JF8ZcyNgjqBTU3C82J0oa+MeAUO5
6ZLLUEaltLJMFbticQkwB1ed8rCSZQDWkkQnIzHOml/wkOJPePR/lX1GAKhrXdZo
I3o+vC8RHPDhHG+26acsl15tN2Fj8uJshVdrAFj6mefiLyJGPG+pSWECgYEA1qmX
mh76c+/P0EMMRgQHCeOQkU/XicyLK05IAl+flUyRCZYqUsnviEmvK1We9agdRMtr
Eg7IjnlqVDOU628hGareaU1WRBPFm0Ov4hIfb67yp36r0uFRir+KBaFxKw6IVvAV
iBv+zYOmJVyI2kx75ZGDButmqzB+YbD6Li2GYpcCgYBEwrNvP6EdA8nJY69Niu1/
P/uxvqymppck7pkRu4j7pFusMvtHeTG7Pu0+nu5LEXlGE8eocjkSyehEbyIX+Ljz
GcGtwVFdymqy4gA4EGTmY/4prSY3UOwr9sEu/lMTbyhCwRYMRg+2Znx54EksFAI3
/yDuvjutRhpxww6qiMn0YQKBgQDE/36Rgjx2iW49wkpRNwD+okjaElvLqltNstmC
1B6v9URld9n/gDLC8FxBeKIY799scwIf4FFN7z8VZwETzzihRZ43JTI1569Bfiy7
W8ZdyEAIVsd5EC61FnKkGDSzPvMAVfRspMSB9n9TakhtjiNl2tRUVVQzZp2VKcVu
+3iIoQKBgQCdihCwCX3kNx+JnR9YDsmvSHeW6A69dSuBdqRAuaWT8Ru+zjOe9a39
IVymlFwUuVhS5pnkPYPLssXNPtaAAC2kJvZE0XVc3wbOlw0XRIadKlHrMZjb9xSj
XrQRqjLVg9YirOmCy390UsEz+FDq8RcE+lPCTUrHDfSXb2/aqQvguw==
-----END RSA PRIVATE KEY-----`
)

var _ = Describe("TLS key rotation", func() {
	var (
		manifest     string
		manifestName string

		kv      consulclient.HTTPKV
		spammer *helpers.Spammer
	)

	BeforeEach(func() {
		var err error
		manifest, err = helpers.DeployConsulWithInstanceCount("tls-key-rotation", 3, config.WindowsClients, boshClient)
		Expect(err).NotTo(HaveOccurred())

		manifestName, err = ops.ManifestName(manifest)
		Expect(err).NotTo(HaveOccurred())

		Eventually(func() ([]bosh.VM, error) {
			return helpers.DeploymentVMs(boshClient, manifestName)
		}, "5m", "10s").Should(ConsistOf(helpers.GetVMsFromManifest(manifest)))

		testConsumerIPs, err := helpers.GetVMIPs(boshClient, manifestName, "testconsumer")
		Expect(err).NotTo(HaveOccurred())

		kv = consulclient.NewHTTPKV(fmt.Sprintf("http://%s:6769", testConsumerIPs[0]))

		spammer = helpers.NewSpammer(kv, 1*time.Second, "testconsumer")
	})

	AfterEach(func() {
		if !CurrentGinkgoTestDescription().Failed {
			err := boshClient.DeleteDeployment(manifestName)
			Expect(err).NotTo(HaveOccurred())
		}
	})

	It("successfully rolls with new tls keys and certs", func() {
		By("spamming the kv store", func() {
			spammer.Spam()
		})

		By("adding a new ca cert", func() {
			oldCACert, err := ops.FindOp(manifest, "/instance_groups/name=consul/properties/consul/ca_cert")
			Expect(err).NotTo(HaveOccurred())

			manifest, err = ops.ApplyOp(manifest, ops.Op{
				Type:  "replace",
				Path:  "/instance_groups/name=consul/properties/consul/ca_cert",
				Value: fmt.Sprintf("%s\n%s", oldCACert.(string), newCACert),
			})
			Expect(err).NotTo(HaveOccurred())
		})

		By("deploying with the new ca cert", func() {
			_, err := boshClient.Deploy([]byte(manifest))
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() ([]bosh.VM, error) {
				return helpers.DeploymentVMs(boshClient, manifestName)
			}, "5m", "10s").Should(ConsistOf(helpers.GetVMsFromManifest(manifest)))
		})

		By("replace agent and server keys and certs", func() {
			var err error
			manifest, err = ops.ApplyOps(manifest, []ops.Op{
				{
					Type:  "replace",
					Path:  "/instance_groups/name=consul/properties/consul/agent_cert",
					Value: newAgentCert,
				},
				{
					Type:  "replace",
					Path:  "/instance_groups/name=consul/properties/consul/server_cert",
					Value: newServerCert,
				},
				{
					Type:  "replace",
					Path:  "/instance_groups/name=consul/properties/consul/agent_key",
					Value: newAgentKey,
				},
				{
					Type:  "replace",
					Path:  "/instance_groups/name=consul/properties/consul/server_key",
					Value: newServerKey,
				},
			})
			Expect(err).NotTo(HaveOccurred())
		})

		By("deploying with the new agent and server keys and certs", func() {
			_, err := boshClient.Deploy([]byte(manifest))
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() ([]bosh.VM, error) {
				return helpers.DeploymentVMs(boshClient, manifestName)
			}, "5m", "10s").Should(ConsistOf(helpers.GetVMsFromManifest(manifest)))
		})

		By("removing the old ca cert", func() {
			var err error
			manifest, err = ops.ApplyOp(manifest, ops.Op{
				Type:  "replace",
				Path:  "/instance_groups/name=consul/properties/consul/ca_cert",
				Value: newCACert,
			})
			Expect(err).NotTo(HaveOccurred())
		})

		By("deploying with the old ca cert removed", func() {
			_, err := boshClient.Deploy([]byte(manifest))
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() ([]bosh.VM, error) {
				return helpers.DeploymentVMs(boshClient, manifestName)
			}, "5m", "10s").Should(ConsistOf(helpers.GetVMsFromManifest(manifest)))
		})

		By("stopping the spammer", func() {
			spammer.Stop()
		})

		By("reading from the consul kv store", func() {
			err := spammer.Check()
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
