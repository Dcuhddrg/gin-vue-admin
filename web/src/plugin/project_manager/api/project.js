import service from '@/utils/request'

export const createProject = (data) => {
  return service({
    url: '/project/createProject',
    method: 'post',
    data
  })
}

export const deleteProject = (data) => {
  return service({
    url: '/project/deleteProject',
    method: 'delete',
    data
  })
}

export const updateProject = (data) => {
  return service({
    url: '/project/updateProject',
    method: 'put',
    data
  })
}

export const getProjectList = (params) => {
  return service({
    url: '/project/getProjectList',
    method: 'get',
    params
  })
}
