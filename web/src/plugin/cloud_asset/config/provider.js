export const PROVIDER_OPTIONS = [
  { label: '阿里云', value: 'aliyun' },
  { label: '腾讯云', value: 'tencent' },
  { label: 'AWS', value: 'aws' },
  { label: '华为云', value: 'huawei' },
  { label: '百度云', value: 'baidu' }
]

export const getProviderLabel = (type) => {
  const option = PROVIDER_OPTIONS.find(item => item.value === type)
  return option ? option.label : type
}
