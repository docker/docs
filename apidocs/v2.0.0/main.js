

$(function () {
  var url = window.location.search.match(/url=([^&]+)/);
  if (url && url.length > 1) {
    url = decodeURIComponent(url[1]);
  } else {
    url = "../docs.json";
  }

  // Pre load translate...
  if(window.SwaggerTranslator) {
    window.SwaggerTranslator.translate();
  }
  window.swaggerUi = new SwaggerUi({
    spec:
{"swagger":"2.0","info":{"title":"DTR 2.0.0 API Documentation","version":"1.0.0","description":"Docker Trusted Registry has an experimental API that you can use to manage\nDTR repositories, permissions, and settings.\n\n**This API is experimental and subject to change, which could affect future\nbackwards compatibility.**\n"},"tags":[{"name":"accounts","description":"Accounts"},{"name":"index","description":"Index"},{"name":"meta","description":"Admin"},{"name":"openid","description":"OpenID"},{"name":"repositories","description":"Repositories"},{"name":"repositoryNamespaces","description":"Repository Namespaces"}],"paths":{"/api/v0/accounts/{orgname}/teams/{teamname}/repositoryAccess":{"get":{"tags":["accounts"],"operationId":"ListTeamRepoAccess","summary":"List repository access grants for a team","description":"\n*Authorization:* Client must be authenticated as a user who owns the organization the team is in or be a member of that team.\n\t\t","produces":["application/json"],"consumes":["application/json"],"parameters":[{"in":"path","description":"organization account name","name":"orgname","required":true,"type":"string"},{"in":"path","description":"team name","name":"teamname","required":true,"type":"string"},{"in":"query","description":"Start of page index (This can have slightly different meanings for different endpoints)","name":"start","required":false,"type":"integer","format":"int32","default":0},{"in":"query","description":"Maximum number of results to return","name":"limit","required":false,"type":"integer","format":"int32","default":10}],"responses":{"200":{"description":"","schema":{"$ref":"#/definitions/responses.ListTeamRepoAccess"}},"400":{"description":"the team does not belong to the organization"},"401":{"description":"NOT_AUTHENTICATED: The client is not authenticated."},"403":{"description":"NOT_AUTHORIZED: The client is not authorized."},"404":{"description":"NO_SUCH_TEAM: A team with the given name does not exist in the organization."}}}},"/api/v0/accounts/{username}/repositoryAccess/{namespace}/{reponame}":{"get":{"tags":["accounts"],"operationId":"GetUserRepoAccess","summary":"Check a user's access to a repository","description":"\n\t*Authorization:* Client must be authenticated either as the user in question or be a system admin.\n\t\t","produces":["application/json"],"consumes":["application/json"],"parameters":[{"in":"path","description":"user account name","name":"username","required":true,"type":"string"},{"in":"path","description":"namespace of repository","name":"namespace","required":true,"type":"string"},{"in":"path","description":"name of repository","name":"reponame","required":true,"type":"string"}],"responses":{"200":{"description":"","schema":{"$ref":"#/definitions/responses.RepoUserAccess"}},"401":{"description":"NOT_AUTHENTICATED: The client is not authenticated."},"403":{"description":"NOT_AUTHORIZED: The client is not authorized."},"404":{"description":"NO_SUCH_REPOSITORY: A repository with the given name does not exist."}}}},"/api/v0/index/dockersearch":{"get":{"tags":["index"],"operationId":"Docker Search","summary":"Search Docker repositories","description":"\nThis is used for the Docker CLI's docker search command. Repository results will be filtered to only those repositories visible to the client.\n\t\t","produces":["application/json"],"consumes":["application/json"],"parameters":[{"in":"query","description":"Search query","name":"query","required":true,"type":"string"}],"responses":{"200":{"description":"","schema":{"$ref":"#/definitions/responses.DockerSearch"}},"401":{"description":"NOT_AUTHENTICATED: The client is not authenticated."}}}},"/api/v0/index/autocomplete":{"get":{"tags":["index"],"operationId":"Autocomplete","summary":"Autocompletion for repositories and/or accounts","description":"\nRepository results will be filtered to only those repositories visible to the client. Account results will not be filtered.\n\t\t","produces":["application/json"],"consumes":["application/json"],"parameters":[{"in":"query","description":"Autocomplete query","name":"query","required":true,"type":"string"},{"in":"query","description":"Whether to include repositories in the response","name":"includeRepositories","required":false,"type":"boolean","default":true},{"in":"query","description":"Whether to include accounts in the response","name":"includeAccounts","required":false,"type":"boolean","default":true},{"in":"query","description":"Exact repository namespace to limit results to.","name":"namespace","required":false,"type":"string"}],"responses":{"200":{"description":"","schema":{"$ref":"#/definitions/responses.Autocomplete"}},"401":{"description":"NOT_AUTHENTICATED: The client is not authenticated."}}}},"/api/v0/meta/settings":{"post":{"tags":["meta"],"operationId":"UpdateSettings","summary":"Update settings","description":"\n*Authorization:* Client must be authenticated an admin.\n\t\t","produces":["application/json"],"consumes":["application/json"],"parameters":[{"in":"body","name":"body","required":true,"schema":{"$ref":"#/definitions/forms.Settings"}}],"responses":{"202":{"description":"success"},"400":{"description":"INVALID_SETTINGS: The submitted settings change request contains invalid values."},"401":{"description":"NOT_AUTHENTICATED: The client is not authenticated."}}},"get":{"tags":["meta"],"operationId":"GetSettings","summary":"Get settings","description":"\n*Authorization:* Client must be authenticated an admin.\n\t\t","produces":["application/json"],"consumes":["application/json"],"responses":{"200":{"description":"","schema":{"$ref":"#/definitions/responses.Settings"}},"401":{"description":"NOT_AUTHENTICATED: The client is not authenticated."}}}},"/api/v0/meta/cluster_status":{"get":{"tags":["meta"],"operationId":"GetClusterStatus","summary":"Get cluster status","description":"\n*Authorization:* Client must be authenticated an admin.\n\t\t","produces":["application/json"],"consumes":["application/json"],"responses":{"200":{"description":"","schema":{"$ref":"#/definitions/responses.ClusterStatus"}},"401":{"description":"NOT_AUTHENTICATED: The client is not authenticated."}}}},"/api/v0/openid/begin":{"get":{"tags":["openid"],"operationId":"OpenIDBegin","produces":["application/json"],"consumes":["application/json"],"responses":{"200":{"description":""},"401":{"description":"NOT_AUTHENTICATED: The client is not authenticated."}}}},"/api/v0/openid/callback":{"get":{"tags":["openid"],"operationId":"OpenIDCallback","produces":["application/json"],"consumes":["application/json"],"responses":{"200":{"description":""},"401":{"description":"NOT_AUTHENTICATED: The client is not authenticated."}}}},"/api/v0/openid/keys":{"get":{"tags":["openid"],"operationId":"OpenIDKeys","produces":["application/json"],"consumes":["application/json"],"responses":{"200":{"description":"","schema":{"$ref":"#/definitions/responses.OpenIDKeys"}},"401":{"description":"NOT_AUTHENTICATED: The client is not authenticated."}}}},"/api/v0/repositories":{"get":{"tags":["repositories"],"operationId":"ListRepositories","summary":"List all repositories","description":"\n*Authorization:* Client must be authenticated as any active user in the system. Results will be filtered to only those repositories visible to the client.\n\t\t","produces":["application/json"],"consumes":["application/json"],"parameters":[{"in":"query","description":"Start of page index (This can have slightly different meanings for different endpoints)","name":"start","required":false,"type":"integer","format":"int32","default":0},{"in":"query","description":"Maximum number of results to return","name":"limit","required":false,"type":"integer","format":"int32","default":10}],"responses":{"200":{"description":"","schema":{"$ref":"#/definitions/responses.Repositories"}},"401":{"description":"NOT_AUTHENTICATED: The client is not authenticated."}}}},"/api/v0/repositories/{namespace}":{"get":{"tags":["repositories"],"operationId":"ListNamespaceRepositories","summary":"List repositories in a namespace","description":"\n*Authorization:* Client must be authenticated as any active user in the system. Results will be filtered to only those repositories visible to the client.\n\t\t","produces":["application/json"],"consumes":["application/json"],"parameters":[{"in":"path","description":"namespace of repository","name":"namespace","required":true,"type":"string"},{"in":"query","description":"Start of page index (This can have slightly different meanings for different endpoints)","name":"start","required":false,"type":"integer","format":"int32","default":0},{"in":"query","description":"Maximum number of results to return","name":"limit","required":false,"type":"integer","format":"int32","default":10}],"responses":{"200":{"description":"","schema":{"$ref":"#/definitions/responses.Repositories"}},"401":{"description":"NOT_AUTHENTICATED: The client is not authenticated."},"404":{"description":"NO_SUCH_ACCOUNT: An account with the given name does not exist."}}},"post":{"tags":["repositories"],"operationId":"CreateRepository","summary":"Create repository","description":"\n*Authorization:* Client must be authenticated as a user who has admin access to the\nrepository namespace (i.e., user owns the repo or is a member of a team with\n\"admin\" level access to the organization's namespace of repositories).\n\t\t","produces":["application/json"],"consumes":["application/json"],"parameters":[{"in":"path","description":"namespace of repository","name":"namespace","required":true,"type":"string"},{"in":"body","name":"body","required":true,"schema":{"$ref":"#/definitions/forms.CreateRepo"}}],"responses":{"201":{"description":"success","schema":{"$ref":"#/definitions/responses.Repository"}},"400":{"description":"REPOSITORY_EXISTS: A repository with the same name already exists."},"401":{"description":"NOT_AUTHENTICATED: The client is not authenticated."},"403":{"description":"NOT_AUTHORIZED: The client is not authorized."},"404":{"description":"NO_SUCH_ACCOUNT: An account with the given name does not exist."}}}},"/api/v0/repositories/{namespace}/{reponame}":{"get":{"tags":["repositories"],"operationId":"GetRepository","summary":"View details of a repository","description":"\n*Authorization:* Client must be authenticated as a user who has visibility to the repository.\n\t\t","produces":["application/json"],"consumes":["application/json"],"parameters":[{"in":"path","description":"namespace of repository","name":"namespace","required":true,"type":"string"},{"in":"path","description":"name of repository","name":"reponame","required":true,"type":"string"}],"responses":{"200":{"description":"","schema":{"$ref":"#/definitions/responses.Repository"}},"401":{"description":"NOT_AUTHENTICATED: The client is not authenticated."},"404":{"description":"NO_SUCH_REPOSITORY: A repository with the given name does not exist."}}},"patch":{"tags":["repositories"],"operationId":"PatchRepository","summary":"Update details of a repository","description":"\n*Authorization:* Client must be authenticated as a user who has \"admin\" access to the repository\n(i.e., user owns the repo or is a member of a team with \"admin\" level access to the organization\"s repository).\n\nNote that a repository cannot be renamed this way.\n\t\t","produces":["application/json"],"consumes":["application/json"],"parameters":[{"in":"path","description":"namespace of repository","name":"namespace","required":true,"type":"string"},{"in":"path","description":"name of repository","name":"reponame","required":true,"type":"string"},{"in":"body","name":"body","required":true,"schema":{"$ref":"#/definitions/forms.UpdateRepo"}}],"responses":{"200":{"description":"","schema":{"$ref":"#/definitions/responses.Repository"}},"400":{"description":"INVALID_REPOSITORY_VISIBILITY: The visibility value of the repository is invalid."},"401":{"description":"NOT_AUTHENTICATED: The client is not authenticated."},"403":{"description":"NOT_AUTHORIZED: The client is not authorized."},"404":{"description":"NO_SUCH_REPOSITORY: A repository with the given name does not exist."}}},"delete":{"tags":["repositories"],"operationId":"DeleteRepository","summary":"Remove a repository","description":"\n*Authorization:* Client must be authenticated as a user who has \"admin\" access to the repository\n(i.e., user owns the repo or is a member of a team with \"admin\" level access to the organization\"s repository).\n\t\t","produces":["application/json"],"consumes":["application/json"],"parameters":[{"in":"path","description":"namespace of repository","name":"namespace","required":true,"type":"string"},{"in":"path","description":"name of repository","name":"reponame","required":true,"type":"string"}],"responses":{"204":{"description":"success or repository does not exist"},"401":{"description":"NOT_AUTHENTICATED: The client is not authenticated."},"403":{"description":"NOT_AUTHORIZED: The client is not authorized."},"404":{"description":"NO_SUCH_REPOSITORY: A repository with the given name does not exist."}}}},"/api/v0/repositories/{namespace}/{reponame}/teamAccess":{"get":{"tags":["repositories"],"operationId":"ListRepoTeamAccess","summary":"List teams granted access to an organization-owned repository","description":"\n*Authorization:* Client must be authenticated as a user who has \"admin\" level access to the repository.\n\t\t","produces":["application/json"],"consumes":["application/json"],"parameters":[{"in":"path","description":"namespace of repository","name":"namespace","required":true,"type":"string"},{"in":"path","description":"name of repository","name":"reponame","required":true,"type":"string"},{"in":"query","description":"Start of page index (This can have slightly different meanings for different endpoints)","name":"start","required":false,"type":"integer","format":"int32","default":0},{"in":"query","description":"Maximum number of results to return","name":"limit","required":false,"type":"integer","format":"int32","default":10}],"responses":{"200":{"description":"","schema":{"$ref":"#/definitions/responses.ListRepoTeamAccess"}},"400":{"description":"the repository is not owned by an organization"},"401":{"description":"NOT_AUTHENTICATED: The client is not authenticated."},"403":{"description":"NOT_AUTHORIZED: The client is not authorized."},"404":{"description":"NO_SUCH_REPOSITORY: A repository with the given name does not exist."}}}},"/api/v0/repositories/{namespace}/{reponame}/teamAccess/{teamname}":{"put":{"tags":["repositories"],"operationId":"GrantRepoTeamAccess","summary":"Set a team's access to an orgnization-owned repository","description":"\n*Authorization:* Client must be authenticated as a user who has \"admin\" level access to the repository.\n\t\t","produces":["application/json"],"consumes":["application/json"],"parameters":[{"in":"path","description":"namespace of repository","name":"namespace","required":true,"type":"string"},{"in":"path","description":"name of repository","name":"reponame","required":true,"type":"string"},{"in":"path","description":"team name","name":"teamname","required":true,"type":"string"},{"in":"body","name":"body","required":true,"schema":{"$ref":"#/definitions/forms.Access"}}],"responses":{"200":{"description":"","schema":{"$ref":"#/definitions/responses.RepoTeamAccess"}},"400":{"description":"the team does not belong to the organization"},"401":{"description":"NOT_AUTHENTICATED: The client is not authenticated."},"403":{"description":"NOT_AUTHORIZED: The client is not authorized."},"404":{"description":"NO_SUCH_TEAM: A team with the given name does not exist in the organization."}}},"delete":{"tags":["repositories"],"operationId":"RevokeRepoTeamAccess","summary":"Revoke a team's acccess to an organization-owned repository","description":"\n*Authorization:* Client must be authenticated as a user who has \"admin\" level access to the repository.\n\t\t","produces":["application/json"],"consumes":["application/json"],"parameters":[{"in":"path","description":"namespace of repository","name":"namespace","required":true,"type":"string"},{"in":"path","description":"name of repository","name":"reponame","required":true,"type":"string"},{"in":"path","description":"team name","name":"teamname","required":true,"type":"string"}],"responses":{"204":{"description":"success or the team is not in the access list or there is no such team in the organization"},"400":{"description":"the repository is not owned by an organization"},"401":{"description":"NOT_AUTHENTICATED: The client is not authenticated."},"403":{"description":"NOT_AUTHORIZED: The client is not authorized."},"404":{"description":"NO_SUCH_TEAM: A team with the given name does not exist in the organization."}}}},"/api/v0/repositories/{namespace}/{reponame}/tags":{"get":{"tags":["repositories"],"operationId":"ListRepoTags","summary":"List the available tags for a repository","description":"\n*Authorization:* Client must be authenticated as a user who has visibility to the repository.\n\t\t","produces":["application/json"],"consumes":["application/json"],"parameters":[{"in":"path","description":"namespace of repository","name":"namespace","required":true,"type":"string"},{"in":"path","description":"name of repository","name":"reponame","required":true,"type":"string"}],"responses":{"200":{"description":"","schema":{"$ref":"#/definitions/responses.ListRepositoryTags"}},"401":{"description":"NOT_AUTHENTICATED: The client is not authenticated."},"404":{"description":"NO_SUCH_REPOSITORY: A repository with the given name does not exist."}}}},"/api/v0/repositories/{namespace}/{reponame}/tags/{tag}/trust":{"get":{"tags":["repositories"],"operationId":"GetTrustForTag","summary":"Get Notary trust info about a specific tag","description":"\nRepository results will be filtered to only those repositories visible to the client.\n\t\t","produces":["application/json"],"consumes":["application/json"],"parameters":[{"in":"path","description":"namespace of repository","name":"namespace","required":true,"type":"string"},{"in":"path","description":"name of repository","name":"reponame","required":true,"type":"string"},{"in":"path","description":"tag name","name":"tag","required":true,"type":"string"}],"responses":{"200":{"description":"","schema":{"$ref":"#/definitions/responses.TagWithSignature"}},"401":{"description":"NOT_AUTHENTICATED: The client is not authenticated."},"404":{"description":"NO_SUCH_TAG: A tag with the given name does not exist for the given repository."}}}},"/api/v0/repositories/{namespace}/{reponame}/manifests/{reference}":{"delete":{"tags":["repositories"],"operationId":"DeleteRepoManifest","summary":"Delete a manifest for a repository","description":"\n*Authorization:* Client must be authenticated as a user who has \"write\" level access to the repository.\n\t\t","produces":["application/json"],"consumes":["application/json"],"parameters":[{"in":"path","description":"namespace of repository","name":"namespace","required":true,"type":"string"},{"in":"path","description":"name of repository","name":"reponame","required":true,"type":"string"},{"in":"path","description":"digest or tag for an image manifest","name":"reference","required":true,"type":"string"}],"responses":{"202":{"description":"success"},"401":{"description":"NOT_AUTHENTICATED: The client is not authenticated."},"403":{"description":"NOT_AUTHORIZED: The client is not authorized."},"404":{"description":"NO_SUCH_MANIFEST: A manifest with the given reference does not exist for the given repository."}}}},"/api/v0/repositoryNamespaces/{namespace}/teamAccess":{"get":{"tags":["repositoryNamespaces"],"operationId":"ListRepoNamespaceTeamAccess","summary":"List teams granted access to an organization-owned namespace of repositories","description":"\n*Authorization:* Client must be authenticated as a user who has ‘admin’ level access to the namespace.\n\t\t","produces":["application/json"],"consumes":["application/json"],"parameters":[{"in":"path","description":"namespace of repository","name":"namespace","required":true,"type":"string"},{"in":"query","description":"Start of page index (This can have slightly different meanings for different endpoints)","name":"start","required":false,"type":"integer","format":"int32","default":0},{"in":"query","description":"Maximum number of results to return","name":"limit","required":false,"type":"integer","format":"int32","default":10}],"responses":{"200":{"description":"","schema":{"$ref":"#/definitions/responses.ListRepoNamespaceTeamAccess"}},"400":{"description":"the namespace is not owned by an organization"},"401":{"description":"NOT_AUTHENTICATED: The client is not authenticated."},"403":{"description":"NOT_AUTHORIZED: The client is not authorized."},"404":{"description":"NO_SUCH_ACCOUNT: An account with the given name does not exist."}}}},"/api/v0/repositoryNamespaces/{namespace}/teamAccess/{teamname}":{"get":{"tags":["repositoryNamespaces"],"operationId":"GetRepoNamespaceTeamAccess","summary":"Get a team's granted access to an organization-owned namespace of repositories","description":"\n*Authorization:* Client must be authenticated as a user who has \"admin\" level access to\nthe namespace, is a system admin, member of the organization's \"owners\" team, or is a\nmember of the team in question.\n\t\t","produces":["application/json"],"consumes":["application/json"],"parameters":[{"in":"path","description":"namespace of repository","name":"namespace","required":true,"type":"string"},{"in":"path","description":"team name","name":"teamname","required":true,"type":"string"}],"responses":{"200":{"description":"","schema":{"$ref":"#/definitions/responses.NamespaceTeamAccess"}},"401":{"description":"NOT_AUTHENTICATED: The client is not authenticated."},"403":{"description":"NOT_AUTHORIZED: The client is not authorized."},"404":{"description":"NO_SUCH_TEAM: A team with the given name does not exist in the organization."}}},"put":{"tags":["repositoryNamespaces"],"operationId":"GrantRepoNamespaceTeamAccess","summary":"Set a team's access to an organization-owned namespace of repositories","description":"\n*Authorization:* Client must be authenticated as a user who has ‘admin’ level access to the namespace.\n\t\t","produces":["application/json"],"consumes":["application/json"],"parameters":[{"in":"path","description":"namespace of repository","name":"namespace","required":true,"type":"string"},{"in":"path","description":"team name","name":"teamname","required":true,"type":"string"},{"in":"body","name":"body","required":true,"schema":{"$ref":"#/definitions/forms.Access"}}],"responses":{"200":{"description":"","schema":{"$ref":"#/definitions/responses.NamespaceTeamAccess"}},"400":{"description":"the team does not belong to the owning organization"},"401":{"description":"NOT_AUTHENTICATED: The client is not authenticated."},"403":{"description":"NOT_AUTHORIZED: The client is not authorized."},"404":{"description":"NO_SUCH_TEAM: A team with the given name does not exist in the organization."}}},"delete":{"tags":["repositoryNamespaces"],"operationId":"RevokeRepoNamespaceTeamAccess","summary":"Revoke a team's access to an organization-owned namespace of repositories","description":"\n*Authorization:* Client must be authenticated as a user who has ‘admin’ level access to the namespace.\n\t\t","produces":["application/json"],"consumes":["application/json"],"parameters":[{"in":"path","description":"namespace of repository","name":"namespace","required":true,"type":"string"},{"in":"path","description":"team name","name":"teamname","required":true,"type":"string"}],"responses":{"204":{"description":"success or the team does not exist in the access list or there is no such team in the organization"},"401":{"description":"NOT_AUTHENTICATED: The client is not authenticated."},"403":{"description":"NOT_AUTHORIZED: The client is not authorized."},"404":{"description":"NO_SUCH_TEAM: A team with the given name does not exist in the organization."}}}}},"definitions":{"responses.ListTeamRepoAccess":{"required":["team","repositoryAccessList"],"properties":{"team":{"$ref":"#/definitions/responses.Team"},"repositoryAccessList":{"type":"array","items":{"$ref":"#/definitions/responses.RepoAccess"}}}},"responses.Team":{"required":["id","clientUserIsMember"],"properties":{"id":{"type":"string"},"clientUserIsMember":{"type":"boolean"}}},"responses.RepoAccess":{"required":["accessLevel","repository"],"properties":{"accessLevel":{"type":"string","enum":["read-only","read-write","admin"]},"repository":{"$ref":"#/definitions/responses.Repository"}}},"responses.Repository":{"required":["id","namespace","namespaceType","name","shortDescription","visibility"],"properties":{"id":{"type":"string"},"namespace":{"type":"string"},"namespaceType":{"type":"string","enum":["user","organization"]},"name":{"type":"string"},"shortDescription":{"type":"string"},"longDescription":{"type":"string"},"visibility":{"type":"string","enum":["public","private"]}}},"responses.RepoUserAccess":{"required":["accessLevel","user","repository"],"properties":{"accessLevel":{"type":"string","enum":["read-only","read-write","admin"]},"user":{"$ref":"#/definitions/responses.Account"},"repository":{"$ref":"#/definitions/responses.Repository"}}},"responses.Account":{"required":["name","id","fullName","isOrg"],"properties":{"name":{"type":"string","description":"Name of the account"},"id":{"type":"string","description":"ID of the account"},"fullName":{"type":"string","description":"Full Name of the account"},"isOrg":{"type":"boolean","description":"Whether the account is an organization (or user)"},"isAdmin":{"type":"boolean","description":"Whether the user is a system admin (users only)"},"isActive":{"type":"boolean","description":"Whether the user is active and can login (users only)"}}},"responses.DockerSearch":{"required":["num_results","query","results"],"properties":{"num_results":{"type":"integer","format":"int32"},"query":{"type":"string"},"results":{"type":"array","items":{"$ref":"#/definitions/responses.DockerRepository"}}}},"responses.DockerRepository":{"required":["description","is_official","is_trusted","name","star_count"],"properties":{"description":{"type":"string"},"is_official":{"type":"boolean"},"is_trusted":{"type":"boolean"},"name":{"type":"string"},"star_count":{"type":"integer","format":"int32"}}},"responses.Autocomplete":{"properties":{"accountResults":{"type":"array","items":{"$ref":"#/definitions/responses.Account"}},"repositoryResults":{"type":"array","items":{"$ref":"#/definitions/responses.Repository"}}}},"forms.Settings":{"properties":{"dtrHost":{"type":"string"},"notaryServer":{"type":"string"},"notaryCert":{"type":"string"},"notaryVerifyCert":{"type":"boolean"},"authBypassCA":{"type":"string"},"authBypassOU":{"type":"string"},"disableUpgrades":{"type":"boolean"},"reportAnalytics":{"type":"boolean"},"releaseChannel":{"type":"string"},"webTLSCert":{"type":"string"},"webTLSKey":{"type":"string"},"webTLSCA":{"type":"string"}}},"responses.Settings":{"required":["dtrHost","replicaSettings","notaryServer","notaryCert","notaryVerifyCert","authBypassCA","authBypassOU","httpProxy","httpsProxy","noProxy","disableUpgrades","reportAnalytics","releaseChannel","logProtocol","logHost","logLevel","webTLSCert","webTLSCA","replicaID"],"properties":{"dtrHost":{"type":"string"},"replicaSettings":{},"notaryServer":{"type":"string"},"notaryCert":{"type":"string"},"notaryVerifyCert":{"type":"boolean"},"authBypassCA":{"type":"string"},"authBypassOU":{"type":"string"},"httpProxy":{"type":"string"},"httpsProxy":{"type":"string"},"noProxy":{"type":"string"},"disableUpgrades":{"type":"boolean"},"reportAnalytics":{"type":"boolean"},"releaseChannel":{"type":"string"},"logProtocol":{"type":"string"},"logHost":{"type":"string"},"logLevel":{"type":"string"},"webTLSCert":{"type":"string"},"webTLSCA":{"type":"string"},"replicaID":{"type":"string"}}},"responses.ClusterStatus":{"required":["rethink_system_tables","etcd_status","replica_health","replica_timestamp","replica_readonly"],"properties":{"rethink_system_tables":{},"etcd_status":{},"replica_health":{},"replica_timestamp":{},"replica_readonly":{}}},"responses.OpenIDKeys":{"required":["keys"],"properties":{"keys":{"type":"array","items":{"$ref":"#/definitions/jose.PublicKey"}}}},"jose.PublicKey":{"required":["kid","kty","publicKey","sigAlgs"],"properties":{"kid":{"type":"string"},"kty":{"type":"string"},"n":{"type":"string"},"e":{"type":"string"},"crv":{"type":"string"},"x":{"type":"string"},"y":{"type":"string"},"publicKey":{"$ref":"#/definitions/crypto.PublicKey"},"sigAlgs":{"type":"array","items":{"$ref":"#/definitions/jose.SignatureAlgorithm"}}}},"jose.SignatureAlgorithm":{},"responses.Repositories":{"required":["repositories"],"properties":{"repositories":{"type":"array","items":{"$ref":"#/definitions/responses.Repository"}}}},"forms.CreateRepo":{"required":["name","shortDescription","longDescription"],"properties":{"name":{"type":"string"},"shortDescription":{"type":"string"},"longDescription":{"type":"string"},"visibility":{"type":"string","enum":["public","private"]}}},"forms.UpdateRepo":{"properties":{"shortDescription":{"type":"string"},"longDescription":{"type":"string"},"visibility":{"type":"string","enum":["public","private"]}}},"responses.ListRepoTeamAccess":{"required":["repository","teamAccessList"],"properties":{"repository":{"$ref":"#/definitions/responses.Repository"},"teamAccessList":{"type":"array","items":{"$ref":"#/definitions/responses.TeamAccess"}}}},"responses.TeamAccess":{"required":["accessLevel","team"],"properties":{"accessLevel":{"type":"string","enum":["read-only","read-write","admin"]},"team":{"$ref":"#/definitions/responses.Team"}}},"forms.Access":{"required":["accessLevel"],"properties":{"accessLevel":{"type":"string","enum":["read-only","read-write","admin"]}}},"responses.RepoTeamAccess":{"required":["accessLevel","team","repository"],"properties":{"accessLevel":{"type":"string","enum":["read-only","read-write","admin"]},"team":{"$ref":"#/definitions/responses.Team"},"repository":{"$ref":"#/definitions/responses.Repository"}}},"responses.ListRepositoryTags":{"required":["name","tags"],"properties":{"name":{"type":"string"},"tags":{"type":"array","items":{"$ref":"#/definitions/responses.TagWithSignature"}}}},"responses.TagWithSignature":{"required":["name","inRegistry","hashMismatch","inNotary"],"properties":{"name":{"type":"string"},"inRegistry":{"type":"boolean","description":"true if the tag exists in Registry"},"hashMismatch":{"type":"boolean","description":"true if the hashes from notary and registry don't match"},"inNotary":{"type":"boolean","description":"true if the tax exists in Notary"}}},"responses.ListRepoNamespaceTeamAccess":{"required":["namespace","teamAccessList"],"properties":{"namespace":{"type":"string"},"teamAccessList":{"type":"array","items":{"$ref":"#/definitions/responses.TeamAccess"}}}},"responses.NamespaceTeamAccess":{"required":["accessLevel","team","namespace"],"properties":{"accessLevel":{"type":"string","enum":["read-only","read-write","admin"]},"team":{"$ref":"#/definitions/responses.Team"},"namespace":{"type":"string"}}}}}
,
    dom_id: "swagger-ui-container",
    supportedSubmitMethods: ['get', 'post', 'put', 'delete', 'patch'],
    validatorUrl: null,
    docExpansion: 'list',
    supportedSubmitMethods: [],
    onComplete: function(swaggerApi, swaggerUi){
      if(typeof initOAuth == "function") {
        initOAuth({
          clientId: "your-client-id",
          clientSecret: "your-client-secret",
          realm: "your-realms",
          appName: "your-app-name", 
          scopeSeparator: ","
        });
      }

      if(window.SwaggerTranslator) {
        window.SwaggerTranslator.translate();
      }

      $('pre code').each(function(i, e) {
        hljs.highlightBlock(e)
      });
    },
    onFailure: function(data) {
      log("Unable to Load SwaggerUI");
    },
    apisSorter: "alpha",
    showRequestHeaders: false
  });

  // if you have an apiKey you would like to pre-populate on the page for demonstration purposes...
  /*
    var apiKey = "myApiKeyXXXX123456789";
    $('#input_apiKey').val(apiKey);
  */

  window.swaggerUi.load();

  function log() {
    if ('console' in window) {
      console.log.apply(console, arguments);
    }
  }
});
