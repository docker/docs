//
//  Options.h
//  Docker
//
//  Created by Doby Mock on 5/16/16.
//  Copyright © 2016 docker. All rights reserved.
//

#ifndef Options_h
#define Options_h

// returns if current user is admin or not
static bool user_info_is_admin() {
    
    struct group *admin_group = getgrnam("admin");
    // if we cannot determine if the user is admin, we consider him as non admin
    if (admin_group == NULL) {
        return false;
    }
    gid_t admin_group_id = admin_group->gr_gid;
    int group_count = getgroups(0, nil);
    gid_t group_list[group_count];
    group_count = getgroups(group_count, group_list);
    for (int i=0; i<group_count; i++) {
        if (group_list[i] == admin_group_id) {
            return true;
        }
    }
    return false;
}

#endif /* Options_h */
