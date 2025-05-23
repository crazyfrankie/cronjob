import request from '@/utils/request'

// 获取任务列表
export function getJobList() {
  return request({
    url: '/api/job',
    method: 'get'
  })
}

// 保存/更新任务
export function saveJob(data) {
  return request({
    url: '/api/job',
    method: 'post',
    data
  })
}

// 删除任务
export function deleteJob(name) {
  return request({
    url: `/api/job/delete`,
    method: 'delete',
    params: { name }
  })
}

// 终止任务
export function killJob(name) {
  return request({
    url: `/api/job/kill`,
    method: 'delete',
    params: { name }
  })
}

// 获取任务日志
export function getJobLogs(name, skip = 0, limit = 20) {
  return request({
    url: '/api/job/logs',
    method: 'get',
    params: { name, skip, limit }
  })
}

// 获取Worker节点列表
export function getWorkerList() {
  return request({
    url: '/api/job/workers',
    method: 'get'
  })
}
