package activitypub

import (
	"reflect"
	"testing"
)

// Reference https://github.com/mastodon/mastodon/blob/main/spec/requests/signature_verification_spec.rb

func TestSignatures(t *testing.T) {

	const privateKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAqIAYvNFGbZ5g4iiK6feSdXD4bDStFM58A7tHycYXaYtzZQpI
eHXAmaXuZzXIwtrP4N0gIk8JNwZvXj2UPS+S07t0V9wNK94he01LV5EMz/GN4eNn
FmDL64HIEuKLvV8TvgjbUPRD6Y5X0UpKi2ZIFLSb96Q5w0Z/k7ntpVKV52y8kz5F
jr/O/0JuHryZe0yItzJh8kzFfeMf0EXzfSnaKvT7P9jhgC6uTre+jXyvVZjiHDrn
qvvucdI3I7DRfXo1OqARBrLjy+TdseUAjNYJ+OuPRI1URIWQI01DCHqcohVu9+Ar
+BiCjFp3ua+XMuJvrvbD61d1Fvig/9nbBRR+8QIDAQABAoIBAAgySHnFWI6gItR3
fkfiqIm80cHCN3Xk1C6iiVu+3oBOZbHpW9R7vl9e/WOA/9O+LPjiSsQOegtWnVvd
RRjrl7Hj20VDlZKv5Mssm6zOGAxksrcVbqwdj+fUJaNJCL0AyyseH0x/IE9T8rDC
I1GH+3tB3JkhkIN/qjipdX5ab8MswEPu8IC4ViTpdBgWYY/xBcAHPw4xuL0tcwzh
FBlf4DqoEVQo8GdK5GAJ2Ny0S4xbXHUURzx/R4y4CCts7niAiLGqd9jmLU1kUTMk
QcXfQYK6l+unLc7wDYAz7sFEHh04M48VjWwiIZJnlCqmQbLda7uhhu8zkF1DqZTu
ulWDGQECgYEA0TIAc8BQBVab979DHEEmMdgqBwxLY3OIAk0b+r50h7VBGWCDPRsC
STD73fQY3lNet/7/jgSGwwAlAJ5PpMXxXiZAE3bUwPmHzgF7pvIOOLhA8O07tHSO
L2mvQe6NPzjZ+6iAO2U9PkClxcvGvPx2OBvisfHqZLmxC9PIVxzruQECgYEAzjM6
BTUXa6T/qHvLFbN699BXsUOGmHBGaLRapFDBfVvgZrwqYQcZpBBhesLdGTGSqwE7
gWsITPIJ+Ldo+38oGYyVys+w/V67q6ud7hgSDTW3hSvm+GboCjk6gzxlt9hQ0t9X
8vfDOYhEXvVUJNv3mYO60ENqQhILO4bQ0zi+VfECgYBb/nUccfG+pzunU0Cb6Dp3
qOuydcGhVmj1OhuXxLFSDG84Tazo7juvHA9mp7VX76mzmDuhpHPuxN2AzB2SBEoE
cSW0aYld413JRfWukLuYTc6hJHIhBTCRwRQFFnae2s1hUdQySm8INT2xIc+fxBXo
zrp+Ljg5Wz90SAnN5TX0AQKBgDaatDOq0o/r+tPYLHiLtfWoE4Dau+rkWJDjqdk3
lXWn/e3WyHY3Vh/vQpEqxzgju45TXjmwaVtPATr+/usSykCxzP0PMPR3wMT+Rm1F
rIoY/odij+CaB7qlWwxj0x/zRbwB7x1lZSp4HnrzBpxYL+JUUwVRxPLIKndSBTza
GvVRAoGBAIVBcNcRQYF4fvZjDKAb4fdBsEuHmycqtRCsnkGOz6ebbEQznSaZ0tZE
+JuouZaGjyp8uPjNGD5D7mIGbyoZ3KyG4mTXNxDAGBso1hrNDKGBOrGaPhZx8LgO
4VXJ+ybXrATf4jr8ccZYsZdFpOphPzz+j55Mqg5vac5P1XjmsGTb
-----END RSA PRIVATE KEY-----`

	t.Run("with an HTTP Signature from a known account", func(t *testing.T) {
		expected := "keyId=\"https://remote.domain/users/bob#main-key\",algorithm=\"rsa-sha256\",headers=\"date host (request-target)\",signature=\"Z8ilar3J7bOwqZkMp7sL8sRs4B1FT+UorbmvWoE+A5UeoOJ3KBcUmbsh+k3wQwbP5gMNUrra9rEWabpasZGphLsbDxfbsWL3Cf0PllAc7c1c7AFEwnewtExI83/qqgEkfWc2z7UDutXc2NfgAx89Ox8DXU/fA2GG0jILjB6UpFyNugkY9rg6oI31UnvfVi3R7sr3/x8Ea3I9thPvqI2byF6cojknSpDAwYzeKdngX3TAQEGzFHz3SDWwyp3jeMWfwvVVbM38FxhvAnSumw7YwWW4L7M7h4M68isLimoT3yfCn2ucBVL5Dz8koBpYf/40w7QidClAwCafZQFC29yDOg==\""

		headers := []Pair{
			{K: "date", V: "Wed, 20 Dec 2023 10:00:00 GMT"},
			{K: "host", V: "www.example.com"},
		}

		result, _ := Sign(privateKey, "https://remote.domain/users/bob#main-key", "get /activitypub/success", headers)
		if !reflect.DeepEqual(expected, result) {
			t.Errorf("result was incorrect\ngot: %v\n\nwant: %v", result, expected)
		}
	})

	t.Run("with a valid signature on a GET request that has a query string", func(t *testing.T) {
		expected := "keyId=\"https://remote.domain/users/bob#main-key\",algorithm=\"rsa-sha256\",headers=\"date host (request-target)\",signature=\"SDMa4r/DQYMXYxVgYO2yEqGWWUXugKjVuz0I8dniQAk+aunzBaF2aPu+4grBfawAshlx1Xytl8lhb0H2MllEz16/tKY7rUrb70MK0w8ohXgpb0qs3YvQgdj4X24L1x2MnkFfKHR/J+7TBlnivq0HZqXm8EIkPWLv+eQxu8fbowLwHIVvRd/3t6FzvcfsE0UZKkoMEX02542MhwSif6cu7Ec/clsY9qgKahb9JVGOGS1op9Lvg/9y1mc8KCgD83U5IxVygYeYXaVQ6gixA9NgZiTCwEWzHM5ELm7w5hpdLFYxYOHg/3G3fiqJzpzNQAcCD4S4JxfE7hMI0IzVlNLT6A==\""

		headers := []Pair{
			{K: "date", V: "Wed, 20 Dec 2023 10:00:00 GMT"},
			{K: "host", V: "www.example.com"},
		}

		result, _ := Sign(privateKey, "https://remote.domain/users/bob#main-key", "get /activitypub/success?foo=42", headers)
		if !reflect.DeepEqual(expected, result) {
			t.Errorf("result was incorrect\ngot: %v\n\nwant: %v", result, expected)
		}
	})

	t.Run("with a valid signature on a POST request", func(t *testing.T) {
		expected := "keyId=\"https://remote.domain/users/bob#main-key\",algorithm=\"rsa-sha256\",headers=\"host date digest (request-target)\",signature=\"gmhMjgMROGElJU3fpehV2acD5kMHeELi8EFP2UPHOdQ54H0r55AxIpji+J3lPe+N2qSb/4H1KXIh6f0lRu8TGSsu12OQmg5hiO8VA9flcA/mh9Lpk+qwlQZIPRqKP9xUEfqD+Z7ti5wPzDKrWAUK/7FIqWgcT/mlqB1R1MGkpMFc/q4CIs2OSNiWgA4K+Kp21oQxzC2kUuYob04gAZ7cyE/FTia5t08uv6lVYFdRsn4XNPn1MsHgFBwBMRG79ng3SyhoG4PrqBEi5q2IdLq3zfre/M6He3wlCpyO2VJNdGVoTIzeZ0Zz8jUscPV3XtWUchpGclLGSaKaq/JyNZeiYQ==\""

		headers := []Pair{
			{K: "host", V: "www.example.com"},
			{K: "date", V: "Wed, 20 Dec 2023 10:00:00 GMT"},
			{K: "digest", V: "Hello world"},
		}

		result, _ := Sign(privateKey, "https://remote.domain/users/bob#main-key", "post /activitypub/success", headers)
		if !reflect.DeepEqual(expected, result) {
			t.Errorf("result was incorrect\ngot: %v\n\nwant: %v", result, expected)
		}
	})
}
