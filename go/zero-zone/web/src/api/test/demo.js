import request from '@/utils/request';

export default {
  time() {
    return request({
      url: '/v1/health',
      method: 'get',
    });
  },
};
