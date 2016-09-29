'use strict';

export default function onAddCollaboratorChange(actionContext, { collaborator }) {
    actionContext.dispatch('ON_ADD_COLLAB_CHANGE', collaborator);
}
