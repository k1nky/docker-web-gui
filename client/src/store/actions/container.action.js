import { request } from '../../utilities/request'
import { toaster } from 'evergreen-ui'
import { logout } from '../../utilities/auth'

export const genericContainer = payload => ({
  type: 'GENERIC_CONTAINER',
  payload
})

export const updateContainer = payload => ({
  type: 'UPDATE_CONTAINER',
  payload
})

export const removeContainer = payload => ({
  type: 'DELETE_CONTAINER',
  payload
})

export const updateContainerLog = payload => ({
  type: 'UPDATE_LOG',
  payload
})

export const toggleModal = payload => ({
  type: 'TOGGLE_MODAL',
  payload
})

export const getContainers = (status = 'running') => {
  return dispatch => {
    dispatch(genericContainer({
      loading: status,
      pageError: false,
      segment: status,
      activeIndex: 0,
      containerListLoading: true,
    }))
    request('get', `container/fetch?status=${status}`, {})
      .then(response => {
        console.log(response)
        //if (response.status == 401) {
        //}
        dispatch(genericContainer({
          loading: false,
          containers: response.data,
          pageError: false,
          containerListLoading: false,
        }))
      }).catch(error => {
        logout()
        window.location.reload(true)
        dispatch(genericContainer({
          loading: false,
          pageError: true,
          containerListLoading: false,
        }))
      })
  }
}

export const toggleContainer = (container, status, hideToaster) => {
  return dispatch => {
    dispatch(updateContainer({
      containerId: container.id,
      data: { stateToggling: true },
    }))
    request('get', `container/command?container=${container.id}&command=${status}`)
      .then(res => {
        const State = {
          ...container.State,
          ...{
            Running: status === 'start'
              ? true
              : false
          }
        }
        dispatch(updateContainer({
          containerId: container.id,
          data: { 
            State,
            stateToggling: false,
          },
        }))
        if(! !!hideToaster) {
          toaster.success(
            `Container ${container.name} has been ${status === 'start'? 'started' : 'stopped'}.`,
            { duration: 5 }
          )
        }
      })
      .catch( ex => {
        dispatch(updateContainer({
          containerId: container.id,
          data: { stateToggling: false },
        }))
        toaster.warning(`There is problem ${status === 'start'? 'starting' : 'stopping'} container ${container.name}`,{duration: 5})
      })
  }
}

export const restartContainer = (container, status) => {
  return dispatch => {
    dispatch(updateContainer({
      containerId: container.id,
      data: { 
        stateToggling: true,
        State: {
          ...container.State,
          ...{
            Running: false
          }
        }
      },
    }))
    request('get', `container/command?container=${container.id}&command=${status}`)
      .then(res => {
        dispatch(updateContainer({
          containerId: container.id,
          data: { 
            State: {
              ...container.running,
              ...{
                Running: true
              }
            },
            stateToggling: false,
          },
        }))
        toaster.success(`Container ${container.name} has been restarted.`,{ duration: 5 })
      })
      .catch( ex => {
        dispatch(updateContainer({
          containerId: container.id,
          data: { 
            State: {
              ...container.State,
              ...{
                Running: false
              }
            },
            stateToggling: false,
          },
        }))
        toaster.warning(`There is problem restarting container ${container.name}`,{duration: 5})
      })
  }
}

export const deleteContainer = (container, command) => (dispatch, getState)=>{
    request('get', `container/command?container=${container.id}&command=${command}`)
      .then(res => {
        dispatch(removeContainer({
          containerId: container.id,
          showModal: !getState().container.showModal,
          selectedContainer: {}
        }))
        toaster.success(
          `Container ${container.name} is no more!!!.`,
          {
            duration: 5
          }
        )
      })  
}

export const getLog = (container) => {
  return dispatch => {
    dispatch(updateContainerLog({
      container: container,
      isShowingSideSheet: false,
    }))
    request('get', `container/logs?container=${container.id}`)
      .then(response => {
        dispatch(updateContainerLog({
          container: container,
          isShowingSideSheet: true,
          logData: response.data
        }))
      })
  }
}

export const resetLogSideSheet = () => (dispatch, getState)=>{
  dispatch(updateContainerLog({
    isShowingSideSheet: !getState().container.isShowingSideSheet,
  }))
}

export const toggleDeleteModal = (container) => (dispatch, getState)=>{
  dispatch(toggleModal({
    showModal: !getState().container.showModal,
    selectedContainer: container ? container : {}
  }))
}