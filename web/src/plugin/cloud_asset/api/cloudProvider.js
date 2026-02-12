import service from '@/utils/request'

export const getCloudProviderList = (data) => {
  return service({
    url: '/cloudProvider/getCloudProviderList',
    method: 'post',
    data: data
  })
}

export const createCloudProvider = (data) => {
  return service({
    url: '/cloudProvider/createCloudProvider',
    method: 'post',
    data: data
  })
}

export const updateCloudProvider = (data) => {
  return service({
    url: '/cloudProvider/updateCloudProvider',
    method: 'put',
    data: data
  })
}

export const deleteCloudProvider = (data) => {
  return service({
    url: '/cloudProvider/deleteCloudProvider',
    method: 'delete',
    data: data
  })
}

export const getCloudProviderById = (data) => {
  return service({
    url: '/cloudProvider/findCloudProvider',
    method: 'get',
    params: data
  })
}

export const getRegions = (data) => {
  return service({
    url: '/cloudProvider/getRegions',
    method: 'post',
    data: data
  })
}

/**
 * 获取厂商配置
 * @param {Object} data
 * @param {number} data.projectId 项目ID
 * @param {string} data.providerType 厂商类型
 */
export const getProviderConfig = (data) => {
  return service({
    url: '/cloudProvider/getProviderConfig',
    method: 'post',
    data: data
  })
}
