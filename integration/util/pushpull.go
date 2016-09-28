package util

import (
	"bytes"
	"encoding/base64"
	"path"
	"strings"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (u *Util) PullTagAndPushImageWithChecks(namespace, reponame, tag, original string) func() {
	err := u.Docker.PullImage(original, nil)
	require.Nil(u.T(), err)
	untag := u.TagImageWithChecks(namespace, reponame, tag, original)
	delTag := u.PushImageWithChecks(namespace, reponame, tag)
	return func() {
		untag()
		delTag()
	}
}

func (u *Util) TagImageWithChecks(namespace, reponame, tag, original string) func() {
	retFunc := func() {}
	err := u.Docker.TagImage(original, path.Join(namespace, reponame), tag, true)
	require.Nil(u.T(), err)
	retFunc = func() {
		_, err := u.Docker.RemoveImage(path.Join(namespace, reponame)+":"+tag, true)
		assert.Nil(u.T(), err)
	}
	return retFunc
}

func (u *Util) PushImageWithChecks(namespace, reponame, tag string) func() {
	retFunc := func() {}
	fullName := path.Join(namespace, reponame)
	err := u.Docker.PushImage(fullName, tag, u.Config.AdminAuthConfig)
	require.Nil(u.T(), err)
	retFunc = func() {
		u.API.Login(u.Config.AdminUsername, u.Config.AdminPassword)
		parts := strings.Split(namespace, "/")
		assert.Nil(u.T(), u.API.DeleteTag(parts[1], reponame, tag))
	}
	return retFunc
}

func (u *Util) LoadPackedImage() func() {
	data, err := base64.StdEncoding.DecodeString(packedImage)
	require.Nil(u.T(), err)

	reader := bytes.NewReader(data)
	err = u.Docker.LoadImage(reader)
	require.Nil(u.T(), err)

	return func() {
		_, err := u.Docker.RemoveImage("dtr/true", true)
		require.Nil(u.T(), err)
	}
}

// packedImage is a base64 encoded .tar.gz of the `dtr/true` image (which is
// tagged from tianon/true) for testing the tagstore
var packedImage = `
H4sIAOWXrFcAA+2a3U4bRxSAnUi96U1fYTW9aSUD8z+7VluRAlUi0VIlNJVCEJqdnYFN7F20u0ZB
yOpz9C36iD1rbANOwaiAI8T5LrzzP7PnzJk5OmuaJUYKFlRmjNCxM5QyoxLBfGAutj6NEx4bz0LI
kjSJjdSpp0E7oZ0LwqZrncVQwCjVPmFsevk5pcM4tGBScaE7lDNYRidStxj7zgzrxlZR1Dmxh5Ut
mmvbLap/pNC76v9DXRaL5gC0lNfpH7Qt5/TPdat/ugwBPHH9n5E8Iz1y121AusRV3ja+HYxTpldo
vMLMLqc9oXqSryZcsoQpod61TcuisXnhK2hs2hGZE3HGhRWCMWeCExaW43RsDHVeSCiDBXIRp4pq
p7PEZiZISXUm5OXhDiAV8kPSOyMvy7op7MDDDNYmMWc2STR30HqzHEDjSR3k/6jH64DUi6ax7uhN
k+UF6QXbr/2lsnLYfF7oq2pWuNucztI7x764Os44t1M4PyvZKk5Ib4/8/mL35Y9rw7pa65fO9tfq
NC96l/Kz7EXFOHGehR+y3yUbg6wdq82u1UfwKivtm377XVEefx9BauPXzWjvPVlrqqF/T/bbLq8G
9rCVQH1kudI9SZVVibVeMG1VMEp4T2Xcapsn1DDurNPKWmZ8kIGBCnSwKsu0C9wrHcMkb8v+cOBr
0iuG/X6X/FlWH/PicDOfSHeraKrT4zIvmmmLneLnYd7Pptltm/o+dD8bjbokK91HUOiJr+q8BDES
tsr4Kj3X9pPWcavDR6BAW7mjvPGuGVZjBQ0y3dpqCfWknxfDT2T0pc8+5B7u/7499dUq3KE3zHHz
/S8/9/+EYhrv/2XQniYPPQdlM/9/7OvPPdvNcFX/DBwIs0z9w1LuOs78yz0S/tra/uX5s2ez/PPO
T50296mzPs6v/0ef9U7cmfZon1/N1V5+Tg/5i8M+Gv/+88Pf31zthyAIgiyTO/t/b7dev3m189uN
c9zs/wFiPv5jJEP/bxmw1UfjqiAPgDKBB5moVHCrhM2CSJROZcykk1yCiYdEKymFCSylQVFrueKx
0GksNYfC1XuI/xqq5+1/yf7/07X/s+sCNRhmexRhNgylYyh90Ra5/YeZ66PuR3ndlBWoZe9swYAx
i+EddazfXXwUOkih51Sm0YqLJvJ8sbkZhbzve4oJCmKjNE2181y7xKfKsEwHpj11jAsTp1bSEINU
oZWDK0pmqYJeUJLyKC+isUoiMuouWuHVb1E3rzCa13eX+MFxc3owjnuSXls82r8S1+6Sqiyb0Iqf
NKfH7c4YN65bAechHORZ3W6xyXZJnE1k4KkC3aeUW5MGDauzqQlxSo3wQXhlvNNeM5lqLZIQYHMx
2FDQgFOyP7pDJH1gizz4urnNPf5/WeT/Cyav/BeAMiGpxPt/GYA5b0xOenIvviBs8tf+uNy1h+NN
njXV2HR6fbCyummPy+1zY4DK+/v4AEaw//WXluVjpAJd1Tkc7bmvH2qORfbPhJi3f6HQ/pfC2cxC
2+tqYqT38H+Q0QjNEUEQBEEQBEEQBEEQBEEQBEEQBEEQBEEQBEEQBEEQ5CH4F2LElMMAUAAA
`
