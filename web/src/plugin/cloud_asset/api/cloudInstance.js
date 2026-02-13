import service from '@/utils/request'

/**
 * 获取云服务器实例列表
 * @param {Object} data 查询参数
 * @param {number} data.page 页码
 * @param {number} data.pageSize 每页数量
 * @param {number} [data.providerId] 云厂商ID
 * @param {string} [data.providerType] 厂商类型
 * @param {string} [data.region] 区域
 * @param {string} [data.instanceId] 实例ID
 * @param {string} [data.instanceName] 实例名称
 * @param {string} [data.status] 实例状态
 * @param {string} [data.publicIp] 公网IP
 * @param {string} [data.privateIp] 内网IP
 * @returns {Promise} 云服务器实例列表
 */
export const getCloudInstanceList = (data) => {
  return service({
    url: '/cloudInstance/getCloudInstanceList',
    method: 'post',
    data: data
  })
}

/**
 * 获取云服务器实例视图列表（带厂商信息）
 * @param {Object} data 查询参数
 * @param {number} data.page 页码
 * @param {number} data.pageSize 每页数量
 * @param {number} [data.providerId] 云厂商ID
 * @param {string} [data.providerType] 厂商类型
 * @param {string} [data.region] 区域
 * @param {string} [data.instanceId] 实例ID
 * @param {string} [data.instanceName] 实例名称
 * @param {string} [data.status] 实例状态
 * @returns {Promise} 云服务器实例视图列表
 */
export const getCloudInstanceVOList = (data) => {
  return service({
    url: '/cloudInstance/getCloudInstanceVOList',
    method: 'post',
    data: data
  })
}

/**
 * 获取云服务器实例详情
 * @param {Object} data 查询参数
 * @param {number} data.id 实例ID
 * @returns {Promise} 云服务器实例详情
 */
export const getCloudInstanceById = (data) => {
  return service({
    url: '/cloudInstance/findCloudInstance',
    method: 'get',
    params: { ID: data.id }
  })
}

/**
 * 创建云服务器实例
 * @param {Object} data 实例数据
 * @returns {Promise} 创建结果
 */
export const createCloudInstance = (data) => {
  return service({
    url: '/cloudInstance/createCloudInstance',
    method: 'post',
    data: data
  })
}

/**
 * 更新云服务器实例
 * @param {Object} data 实例数据
 * @returns {Promise} 更新结果
 */
export const updateCloudInstance = (data) => {
  return service({
    url: '/cloudInstance/updateCloudInstance',
    method: 'put',
    data: data
  })
}

/**
 * 删除云服务器实例
 * @param {Object} data 删除参数
 * @param {number} data.id 实例ID
 * @returns {Promise} 删除结果
 */
export const deleteCloudInstance = (data) => {
  return service({
    url: '/cloudInstance/deleteCloudInstance',
    method: 'delete',
    data: data
  })
}

/**
 * 批量删除云服务器实例
 * @param {Object} data 批量操作参数
 * @param {Array<number>} data.instanceIds 实例ID数组
 * @param {string} data.operation 操作类型（delete）
 * @returns {Promise} 批量删除结果
 */
export const batchDeleteCloudInstance = (data) => {
  return service({
    url: '/cloudInstance/batchDeleteCloudInstance',
    method: 'delete',
    data: data
  })
}

/**
 * 同步云服务器实例
 * @param {Object} data 同步参数
 * @param {number} data.providerId 云厂商ID
 * @param {string} data.region 区域
 * @param {boolean} [data.forceSync] 是否强制同步
 * @returns {Promise} 同步结果
 */
export const syncInstances = (data) => {
  return service({
    url: '/cloudInstance/syncInstances',
    method: 'post',
    data: data
  })
}

/**
 * 批量同步云服务器实例
 * @param {Object} data 批量同步参数
 * @param {number} data.providerId 云厂商ID
 * @param {Array<string>} data.regions 区域数组
 * @param {boolean} [data.forceSync] 是否强制同步
 * @returns {Promise} 批量同步结果
 */
export const batchSyncInstances = (data) => {
  return service({
    url: '/cloudInstance/batchSyncInstances',
    method: 'post',
    data: data
  })
}

/**
 * 获取实例统计信息
 * @returns {Promise} 实例统计数据
 */
export const getInstanceStats = () => {
  return service({
    url: '/cloudInstance/getInstanceStats',
    method: 'get'
  })
}

/**
 * 清除云服务器实例缓存
 * @param {Object} data 清除缓存参数
 * @param {number} [data.providerId] 云厂商ID
 * @param {string} [data.region] 区域
 * @returns {Promise} 清除结果
 */
export const clearInstanceCache = (data) => {
  return service({
    url: '/cloudInstance/clearCache',
    method: 'post',
    data: data
  })
}
