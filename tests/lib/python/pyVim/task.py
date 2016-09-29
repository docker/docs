# VMware vSphere Python SDK
# Copyright (c) 2016 VMware, Inc. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

## @file task.py
## @brief Task functions
##
## This module provies synchronization of client/server operations
## since many VIM operations return 'tasks' which can have
## varying completion times.

"""
Task functions

This module provies synchronization of client/server operations since
many VIM operations return 'tasks' which can have varying completion
times.
"""

from pyVmomi import Vmodl, Vim


##
## @brief Exception class to represent when task is blocked (e.g.:
## waiting for an answer to a question.
##
class TaskBlocked(Exception):
    """
    Exception class to represent when task is blocked (e.g.: waiting
    for an answer to a question.
    """
    pass


#
# TaskUpdates
#     verbose information about task progress
#
def TaskUpdatesVerbose(task, progress):
    if isinstance(task.info.progress, int):
        info = task.info
        if not isinstance(progress, str):
            progress = '%d%% (%s)' % (info.progress, info.state)
        print('Task %s (key:%s, desc:%s) - %s' % (
            info.name.info.name, info.key, info.description, progress))


globalTaskUpdate = None


def SetTasksVerbose(verbose=True):
    global globalTaskUpdate
    if verbose:
        globalTaskUpdate = TaskUpdatesVerbose
    else:
        globalTaskUpdate = None


##
## @param raiseOnError [in] Any exception thrown is thrown up to the caller if
## raiseOnError is set to true
## @param si [in] ServiceInstance to use. If set to None, use the default one.
## @param pc [in] property collector to use else retrieve one from cache
## @param onProgressUpdate [in] callable to call with task progress updates.
##    For example:
##
##    def OnTaskProgressUpdate(task, percentDone):
##       sys.stderr.write('# Task %s: %d%% complete ...\n' % (task, percentDone))
##
## Given a task object and a service instance, wait for the task completion
##
## @return state as either "success" or "error". To look at any errors, the
## user should reexamine the task object.
##
## NOTE: This is a blocking call.
##
def WaitForTask(task,
                raiseOnError=True,
                si=None,
                pc=None,
                onProgressUpdate=None):
    """
    Wait for task to complete.

    @type  raiseOnError      : bool
    @param raiseOnError      : Any exception thrown is thrown up to the caller
                               if raiseOnError is set to true.
    @type  si                : ManagedObjectReference to a ServiceInstance.
    @param si                : ServiceInstance to use. If None, use the
                               information from the task.
    @type  pc                : ManagedObjectReference to a PropertyCollector.
    @param pc                : Property collector to use. If None, get it from
                               the ServiceInstance.
    @type  onProgressUpdate  : callable
    @param onProgressUpdate  : Callable to call with task progress updates.

        For example::

            def OnTaskProgressUpdate(task, percentDone):
                print 'Task %s is %d%% complete.' % (task, percentDone)
    """

    if si is None:
        si = Vim.ServiceInstance("ServiceInstance", task._stub)
    if pc is None:
        pc = si.content.propertyCollector

    progressUpdater = ProgressUpdater(task, onProgressUpdate)
    progressUpdater.Update('created')

    filter = CreateFilter(pc, task)

    version, state = None, None
    # Loop looking for updates till the state moves to a completed state.
    while state not in (Vim.TaskInfo.State.success, Vim.TaskInfo.State.error):
        try:
            version, state = GetTaskStatus(task, version, pc)
            progressUpdater.UpdateIfNeeded()
        except Vmodl.Fault.ManagedObjectNotFound as e:
            print("Task object has been deleted: %s" % e.obj)
            break

    filter.Destroy()

    if state == "error":
        progressUpdater.Update('error: %s' % str(task.info.error))
        if raiseOnError:
            raise task.info.error
        else:
            print("Task reported error: " + str(task.info.error))
    else:
        progressUpdater.Update('completed')

    return state


## Wait for multiple tasks to complete
#  See WaitForTask for detail
#
#  Difference: WaitForTasks won't return the state of tasks. User can check
#  tasks state directly with task.info.state
#
#  TODO: Did not check for question pending
def WaitForTasks(tasks,
                 raiseOnError=True,
                 si=None,
                 pc=None,
                 onProgressUpdate=None,
                 results=None):
    """
    Wait for mulitiple tasks to complete. Much faster than calling WaitForTask
    N times
    """

    if not tasks:
        return

    if si is None:
        si = Vim.ServiceInstance("ServiceInstance", tasks[0]._stub)
    if pc is None:
        pc = si.content.propertyCollector
    if results is None:
        results = []

    progressUpdaters = {}
    for task in tasks:
        progressUpdater = ProgressUpdater(task, onProgressUpdate)
        progressUpdater.Update('created')
        progressUpdaters[str(task)] = progressUpdater

    filter = CreateTasksFilter(pc, tasks)

    try:
        version, state = None, None

        # Loop looking for updates till the state moves to a completed state.
        while len(progressUpdaters):
            update = pc.WaitForUpdates(version)
            for filterSet in update.filterSet:
                for objSet in filterSet.objectSet:
                    task = objSet.obj
                    taskId = str(task)
                    for change in objSet.changeSet:
                        if change.name == 'info':
                            state = change.val.state
                        elif change.name == 'info.state':
                            state = change.val
                        else:
                            continue

                        progressUpdater = progressUpdaters.get(taskId)
                        if not progressUpdater:
                            continue

                        if state == Vim.TaskInfo.State.success:
                            progressUpdater.Update('completed')
                            progressUpdaters.pop(taskId)
                            # cache the results, as task objects could expire if one
                            # of the tasks take a longer time to complete
                            results.append(task.info.result)
                        elif state == Vim.TaskInfo.State.error:
                            err = task.info.error
                            progressUpdater.Update('error: %s' % str(err))
                            if raiseOnError:
                                raise err
                            else:
                                print("Task %s reported error: %s" % (taskId, str(err)))
                                progressUpdaters.pop(taskId)
                        else:
                            if onProgressUpdate:
                                progressUpdater.UpdateIfNeeded()
            # Move to next version
            version = update.version
    finally:
        if filter:
            filter.Destroy()
    return


def GetTaskStatus(task, version, pc):
    update = pc.WaitForUpdates(version)
    state = task.info.state

    if (state == 'running' and task.info.name is not None and task.info.name.info.name != "Destroy"
        and task.info.name.info.name != "Relocate"):
        CheckForQuestionPending(task)

    return update.version, state


def CreateFilter(pc, task):
    """ Create property collector filter for task """
    return CreateTasksFilter(pc, [task])


def CreateTasksFilter(pc, tasks):
    """ Create property collector filter for tasks """
    if not tasks:
        return None

    # First create the object specification as the task object.
    objspecs = [Vmodl.Query.PropertyCollector.ObjectSpec(obj=task)
                for task in tasks]

    # Next, create the property specification as the state.
    propspec = Vmodl.Query.PropertyCollector.PropertySpec(
        type=Vim.Task, pathSet=[], all=True)

    # Create a filter spec with the specified object and property spec.
    filterspec = Vmodl.Query.PropertyCollector.FilterSpec()
    filterspec.objectSet = objspecs
    filterspec.propSet = [propspec]

    # Create the filter
    return pc.CreateFilter(filterspec, True)


def CheckForQuestionPending(task):
    """
    Check to see if VM needs to ask a question, throw exception
    """

    vm = task.info.entity
    if vm is not None and isinstance(vm, Vim.VirtualMachine):
        qst = vm.runtime.question
        if qst is not None:
            raise TaskBlocked("Task blocked, User Intervention required")


##
## @brief Class that keeps track of task percentage complete and calls
## a provided callback when it changes.
##
class ProgressUpdater(object):
    """
    Class that keeps track of task percentage complete and calls a
    provided callback when it changes.
    """

    def __init__(self, task, onProgressUpdate):
        self.task = task
        self.onProgressUpdate = onProgressUpdate
        self.prevProgress = 0
        self.progress = 0

    def Update(self, state):
        global globalTaskUpdate
        taskUpdate = globalTaskUpdate
        if self.onProgressUpdate:
            taskUpdate = self.onProgressUpdate
        if taskUpdate:
            taskUpdate(self.task, state)

    def UpdateIfNeeded(self):
        self.progress = self.task.info.progress

        if self.progress != self.prevProgress:
            self.Update(self.progress)

        self.prevProgress = self.progress
