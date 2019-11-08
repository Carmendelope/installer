/*
 * Copyright 2019 Nalej
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// SCP & SSH Integration tests

package entities

import (
	"encoding/json"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

const privateKeyExample = `
-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAoMkQJnJlmMAJOTnfty+V2Q3SV61ZyhH/BM5Y4Sju9pZUVxOT
lEWMqqiaX2fpBGg+kjAKwCkczUOuwEqHDnw7DqbHskx4tZ+jdxvOxVljJw53kQTD
L6qnZV8Sm9VO8MSD6y0kS+A20UmFgb9QO799Km+KRgebMknnloL2raih+41k5C3+
R5y8rJkaZEL8d9F95qC0uYtWxW+ksbFKxnWPE8h1rWkI0nE3U/xEpFs5vrFuIBDG
EDx5ZLzUjXEGeIt5cEutD19cl4KULMGojQpQWYxR3gGLz/+R5WX89JFYY+1FiBhX
dLmWUYmkmiE3zhRliSaXLm8VUA2n7x91vcgEYQIDAQABAoIBAQCHHvG3ncPLzvbx
ZLWhmRfxRTBUBpbCvsT2IQOIlYHdMRDH7OfFX32Lng29x/GHhqOu7zjZrLNLvWmK
qBdAER8AfSCtsp5u5C3X10K5jxlIpVvOP5ZY5K2w/2kAFQ82P7AtX851BYSL9aGB
HGotDAwAMaSenZ3LcVhyoLT11BXs892T8EhvwRcrCanzJZ6yqafkemoHGVuKGuM/
nHVdMPoasSzPkCnU9snzOLAuGy+oRcLMQwZwMEMXmzAK9OmqYCijP/pr5mujU8je
tIj5dwotCiA20VJPB3R0vaiYEXoM7Ir3MjnXczINpOG8WE1XAKlGOdO1Tm7ma8Zu
O7HDO8zRAoGBAM92ACKZjY+AWLzEMFbZ1tIlvkFCxdq0Zuu/Cfla0SLS9Ag6g/Up
U9FGznDLJjS26c3YWn6Eoc1j9hRSsVfbx+DMn7mmPedMm/czP/Dp5TLcI0yXzSVs
OkUkDJuFl3RyeUt0ROorMqK1WiQNodN67N9HR4g6uGciSUsBgb7u3aRDAoGBAMZn
Zx8KU6OWfCxVM/ABHNnCjBPHq5KyM3qYezsQi/odkC1DnA6/N+qkzp5AfchPUS17
Dlbsyw2KM5MKLkRvF6hQyJNTQxNW1wNe2Qw5TEEJ0b8Lx6p0iG4vBDPDJz8JBCmC
EuWXaAG3yCF6zoOB9gLoEuKhrG9ZzNsWZesFXZyLAoGAQv+7yXDHq9lqTwQZDGNr
ohB4YgEbfqcWOfpHUVVIBzQThXjIVuuS2xo/32NsIkgUN9swVn2k93zZ4vRVu6cJ
5QqQZtdOVJ2EHRBbDQWsdIFtkPXRVc2e/+dFfxBkukGh9IFHJEzxHGTvCIeyhGbF
itItQsyb8wq6mtOQwEXKJJsCgYEAwZiyghJkjLLRlbzKAj5DtaTlZIOoQmuKWe0i
Cf9aZwOj5NcdFzK1UEvipX7OfcAPuS5jTqSeeibJrof3n6U7U20IWuGbCOrqwYoy
hn/+jVQUi7Pl78joO4O7OPsLd7HHku0unUOBJHP9X9XiX2ZX9HwZuXUCumDIyVtw
tcS1lIMCgYAqBSf2j/9gh7EjBQBbXdLLQtuCxz5XlQUVt4iLlbcsZVEmgRSiaqoO
887PH0aaaMUvPvgSNH7tscKAUrQuHWmT9SMCQ4raREx46R5+4OqDie8KMPXKyR7K
JqoSlfnrw6KSKFm1Cu1WqDkRNM6eIHxWw9+HOptTnRDiJRrzikkmuw==
-----END RSA PRIVATE KEY-----
`

var _ = ginkgo.Context("Credentials structure", func() {

	ginkgo.It("Must be able to serialize username and password", func() {
		c := NewCredentials("username", "password")
		serialized, err := json.Marshal(c)
		gomega.Expect(err).To(gomega.BeNil())
		gomega.Expect(serialized).ToNot(gomega.BeNil())
		gomega.Expect(string(serialized)).To(gomega.Equal("{\"username\":\"username\",\"password\":\"password\",\"privateKey\":\"\"}"))
		deserialized := &Credentials{}
		err = json.Unmarshal(serialized, deserialized)
		gomega.Expect(err).To(gomega.BeNil())
		gomega.Expect(deserialized).To(gomega.Equal(c))
	})

	ginkgo.It("Must be able to serialize PKI", func() {
		c := NewPKICredentials("username", privateKeyExample)
		serialized, err := json.Marshal(c)
		gomega.Expect(err).To(gomega.BeNil())
		gomega.Expect(serialized).ToNot(gomega.BeNil())
		deserialized := &Credentials{}
		err = json.Unmarshal(serialized, deserialized)
		gomega.Expect(err).To(gomega.BeNil())
		gomega.Expect(deserialized).To(gomega.Equal(c))
	})
})
