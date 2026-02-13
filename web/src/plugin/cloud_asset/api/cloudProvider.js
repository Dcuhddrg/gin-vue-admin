import service from '@/utils/request'

/**
 * 获取云厂商列表
 * @param {Object} data 查询参数
 * @param {number} data.page 页码
 * @param {number} data.pageSize 每页数量
 * @param {string} [data.name] 厂商名称
 * @param {string} [data.type] 厂商类型
 * @param {number} [data.status] 状态
 * @returns {Promise} 列表数据
 */
export const getCloudProviderList = (data) => {
  return service({
    url: '/cloudProvider/getCloudProviderList',
    method: 'post',
    data: data
  })
}

/**
 * 创建云厂商
 * @param {Object} data 厂商信息
 * @param {string} data.name 名称
 * @param {string} data.type 类型
 * @param {string} data.ak AccessKey
 * @param {string} data.sk SecretKey
 * @param {string} data.region 默认区域
 * @param {number} data.projectId 项目ID
 * @returns {Promise} 创建结果
 */
export const createCloudProvider = (data) => {
  return service({
    url: '/cloudProvider/createCloudProvider',
    method: 'post',
    data: data
  })
}

/**
 * 更新云厂商
 * @param {Object} data 厂商信息
 * @param {number} data.ID ID
 * @param {string} [data.name] 名称
 * @param {string} [data.ak] AccessKey
 * @param {string} [data.sk] SecretKey
 * @returns {Promise} 更新结果
 */
export const updateCloudProvider = (data) => {
  return service({
    url: '/cloudProvider/updateCloudProvider',
    method: 'put',
    data: data
  })
}

/**
 * 删除云厂商
 * @param {Object} data
 * @param {number} data.ID ID
 * @returns {Promise} 删除结果
 */
export const deleteCloudProvider = (data) => {
  return service({
    url: '/cloudProvider/deleteCloudProvider',
    method: 'delete',
    data: data
  })
}

/**
 * 根据ID获取云厂商详情
 * @param {Object} data
 * @param {number} data.ID ID
 * @returns {Promise} 详情数据
 */
export const getCloudProviderById = (data) => {
  return service({
    url: '/cloudProvider/findCloudProvider',
    method: 'get',
    params: data
  })
}

/**
 * 获取云厂商可用区域
 * @param {Object} data
 * @param {number} data.projectId 项目ID
 * @param {string} data.providerType 厂商类型
 * @returns {Promise} 区域列表
 */
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
 * @returns {Promise} 配置信息
 */
export const getProviderConfig = (data) => {
  return service({
    url: '/cloudProvider/getProviderConfig',
    method: 'post',
    data: data
  })
}
