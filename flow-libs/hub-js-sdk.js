type SuperAgentCallback = (err: any,
                           res: any) => void;

type JWT = String;
type ChangePasswordData = {
  username: String;
  oldpassword: String;
  newpassword: String
}

declare module 'hub-js-sdk' {
  declare var Auth: {
    getToken(username: string,
             password: string,
             cb: SuperAgentCallback): void;
  }
  declare var Repositories: {
    createRepository(jwt: JWT,
                     repository: any,
                     cb: SuperAgentCallback): void;
        getReposForUser(jwt: JWT,
                     username: String,
                     cb: SuperAgentCallback): void
  }
  declare var Emails: {
    getEmailSubscriptions(JWT: JWT,
                          user: String,
                          cb: SuperAgentCallback): void;
    unsubscribeEmails(JWT: JWT,
                      user: String,
                      data: Object,
                      cb: SuperAgentCallback): void;
    subscribeEmails(JWT:JWT,
                    user: String,
                    data: Object,
                    cb: SuperAgentCallback): void;
    getEmailsJWT(JWT:JWT,
                 cb:SuperAgentCallback): void;
    getEmailsForUser(JWT: JWT,
                     user: String,
                     cb: SuperAgentCallback): void;
    deleteEmailByID(JWT: JWT,
                    id: String,
                    cb: SuperAgentCallback): void;
    updateEmailByID(JWT: JWT,
                    id: String,
                    data: Object,
                    cb: SuperAgentCallback): void;
    addEmailsForUser(JWT: JWT,
                    user: Object,
                    email: string,
                    cb: SuperAgentCallback): void;
  }
}

declare module 'hub-js-sdk/src/Hub/SDK/Users' {
  declare function changePassword(JWT: JWT,
                                  data: ChangePasswordData,
                                  cb: SuperAgentCallback): void;
  declare function getUser(JWT: JWT,
                           user: String,
                           cb: SuperAgentCallback): void;
}

declare module 'hub-js-sdk/src/Hub/SDK/Auth' {
  declare function getToken(username: String,
                            password: String,
                            cb: SuperAgentCallback): void;
}

declare module 'hub-js-sdk/src/Hub/SDK/Notifications' {
  declare function getActivityFeed(JWT: JWT,
                                  cb: SuperAgentCallback): void;
}
