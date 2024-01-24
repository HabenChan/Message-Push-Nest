<template>
  <div class="main-center-body">
    <div class="container">
      <el-row :gutter="16">
        <el-col :span="8">
          <div class="statistic-card">
            <el-statistic :value="data.today_total_num">
              <template #title>
                <div style="display: inline-flex; align-items: center">
                  今日发送总数
                </div>
              </template>
            </el-statistic>
          </div>
        </el-col>
        <el-col :span="8">
          <div class="statistic-card">
            <el-statistic :value="data.today_succ_num">
              <template #title>
                <div style="display: inline-flex; align-items: center">
                  今日发送成功数
                </div>
              </template>
            </el-statistic>
          </div>
        </el-col>
        <el-col :span="8">
          <div class="statistic-card">
            <el-statistic :value="data.today_failed_num">
              <template #title>
                <div style="display: inline-flex; align-items: center">
                  今日发送失败数
                </div>
              </template>
            </el-statistic>

          </div>
        </el-col>
      </el-row>

      <el-divider />

      <div class="echarts-box">
        <div id="daily-chart" style="width: 65%; height: 350px;"></div>
        <div id="send-cate-chart" style="width: 35%; height: 300px;"></div>
      </div>

    </div>
  </div>
</template>
  
<script >
// import {
//   ArrowRight,
//   CaretBottom,
//   CaretTop,
//   Warning,
// } from '@element-plus/icons-vue'
import { reactive, toRefs, onMounted, onUnmounted } from 'vue'
import { request } from '@/api/api'
import * as echarts from "echarts"

export default {
  setup() {
    const state = reactive({
      data: {},
      dailyChart: {},
      sendCateChart: {},
    });
    let echart = echarts;

    const getStatisticData = async () => {
      const rsp = await request.get('/statistic');
      if (await rsp.data.code == 200) {
        let data = await rsp.data.data;
        state.data = data;
      }
    }

    onMounted(async () => {
      await getStatisticData();
      initDailyChart();
      initSendCateChart();
    });

    onUnmounted(() => {
      state.dailyChart.dispose();
      state.sendCateChart.dispose();
    });

    // 最近30天数据图
    function initDailyChart() {
      state.dailyChart = echart.init(document.getElementById("daily-chart"));

      let xAxisdata = [];
      state.data.latest_send_data.forEach(element => {
        xAxisdata.push(element.day);
      });
      let yAxisdata = [];
      state.data.latest_send_data.forEach(element => {
        yAxisdata.push(element.num);
      });

      state.dailyChart.setOption({
        title: {
          subtext: '最近消息30天发送数据',
          top: 0,
          textStyle: {
            color: '#333',
            fontSize: 18,
          },
        },
        xAxis: {
          type: "category",
          data: xAxisdata
        },
        tooltip: {
          trigger: "axis"
        },
        yAxis: {
          type: "value"
        },
        series: [
          {
            data: yAxisdata,
            type: "line",
            smooth: true
          }
        ]
      });

      window.onresize = function () {
        state.dailyChart.resize();
      };
    }

    // 发送消息实例类别图
    function initSendCateChart() {
      state.sendCateChart = echart.init(document.getElementById("send-cate-chart"));
      let data = [];
      state.data.way_cate_data.forEach(element => {
        data.push({ name: element.way_name, value: element.count_num });
      });
      state.sendCateChart.setOption({
        grid: {
          width: '60%',
          height: '60%',
        },
        title: {
          subtext: '发送消息实例渠道',
          top: 0,
          textStyle: {
            color: '#333',
            fontSize: 18,
          },
        },
        series: [
          {
            type: 'pie',
            data: data,
            radius: '40%',
            // roseType: 'area'
          }
        ]
      });

      window.onresize = function () {
        state.sendCateChart.resize();
      };
    }

    return {
      ...toRefs(state)
    };
  }
}
</script>
  
<style scoped>
.container {
  max-width: 1000px;
  height: 450px;
}


.el-statistic {
  --el-statistic-content-font-size: 28px;
  text-align: center;
}

.statistic-card {
  height: 100%;
  padding: 20px;
  border-radius: 4px;
  background-color: var(--el-bg-color-overlay);
}

.echarts-box {
  display: flex;
}

.green {
  color: var(--el-color-success);
}

.red {
  color: var(--el-color-error);
}
</style>
  