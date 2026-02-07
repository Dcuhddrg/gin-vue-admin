import { mount, config } from '@vue/test-utils'
import { describe, it, expect, vi } from 'vitest'
import RegionsSelector from '../view/components/RegionsSelector.vue'

// Mock Element Plus components
config.global.stubs = {
  'el-icon': true,
  'el-checkbox': {
    template: '<div class="el-checkbox"><slot /></div>',
    props: ['label']
  },
  'el-checkbox-group': {
    template: '<div><slot /></div>',
    props: ['modelValue']
  },
  'el-button': {
    template: '<button><slot /></button>',
    props: ['disabled']
  },
  'InfoFilled': true,
  'Refresh': true
}

// Mock getRegions API
const mockGetRegions = vi.fn()
vi.mock('@/plugin/cloud_asset/api/cloudProvider', () => ({
  getRegions: (...args) => mockGetRegions(...args)
}))

describe('RegionsSelector.vue', () => {
  const defaultProps = {
    provider: 'aliyun',
    accessKey: 'test-ak',
    secretKey: 'test-sk',
    modelValue: ''
  }

  it('renders correctly', () => {
    const wrapper = mount(RegionsSelector, { props: defaultProps })
    expect(wrapper.text()).toContain('AccessKey 与 SecretKey 仅做一次性验证')
    expect(wrapper.text()).toContain('验证并获取可用区')
  })

  it('button disabled when props missing', async () => {
    const wrapper = mount(RegionsSelector, { 
      props: { ...defaultProps, accessKey: '' } 
    })
    // Check computed property
    expect(wrapper.vm.canFetchRegions).toBe('')
  })

  it('calls API and updates regions on click', async () => {
    mockGetRegions.mockResolvedValue({
      code: 0,
      data: [{ regionId: 'cn-hangzhou', localName: '华东1' }]
    })

    const wrapper = mount(RegionsSelector, { props: defaultProps })
    // The button is now inside .flex.justify-end > .el-button
    await wrapper.find('button').trigger('click')
    
    expect(mockGetRegions).toHaveBeenCalledWith({
      provider: 'aliyun',
      accessKey: 'test-ak',
      secretKey: 'test-sk'
    })
    // Check if region list is updated (in stub)
    expect(wrapper.vm.regionState.regions).toEqual([{ regionId: 'cn-hangzhou', localName: '华东1' }])
  })

  it('displays error message on failure', async () => {
    mockGetRegions.mockResolvedValue({
      code: 7,
      msg: 'Invalid Key'
    })

    const wrapper = mount(RegionsSelector, { props: defaultProps })
    await wrapper.find('button').trigger('click')
    
    expect(wrapper.text()).toContain('Invalid Key')
  })

  it('emits update:modelValue on selection', async () => {
    // Setup with data
    mockGetRegions.mockResolvedValue({
      code: 0,
      data: [
        { regionId: 'cn-hangzhou', localName: '杭州' },
        { regionId: 'cn-beijing', localName: '北京' }
      ]
    })
    const wrapper = mount(RegionsSelector, { props: defaultProps })
    
    // Simulate selection change directly
    wrapper.vm.handleSelectionChange(['cn-hangzhou', 'cn-beijing'])
    
    expect(wrapper.emitted('update:modelValue')[0]).toEqual(['cn-hangzhou,cn-beijing'])
  })
})
