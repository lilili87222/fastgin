import request from '@/utils/request'

export function login(data) {
  return request({
    url: '/api/public/login',
    method: 'post',
    data
  })
}

