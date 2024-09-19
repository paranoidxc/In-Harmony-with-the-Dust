<script setup>
import {reactive} from 'vue'
import {Greet, SelectOld, SelectOldFolder, SelectNew, MessageBox, CallCompare} from '../../wailsjs/go/main/App'

const data = reactive({
  name: "",
  resultTip: "",
  compareType: false,
  picked: "false",
  compareDisabled: false,
  old: "",
  new: "",
})

function selectOld() {
  if (data.picked == "false") {
    data.compareType = false
  } else {
    data.compareType = true
  }
  SelectOld(data.compareType).then(result => {
    if (result.length) {
      data.old = result
    }
  })
}

function selectNew() {
  if (data.picked == "false") {
    data.compareType = false
  } else {
    data.compareType = true
  }

  SelectNew(data.compareType).then(result => {
    if (result.length) {
      data.new = result
    }
  })
}

function compare() {
  if (!(data.old.length && data.new.length)) {
    MessageBox("请提供比对文件")
    return
  }

  data.compareDisabled = true
  CallCompare(data.old, data.new).then(result => {
    data.compareDisabled = false
    if (result != "-") {
      data.result = result
    }
    if (result.length == 0) {
      data.resultTip = "文件内容相同"
    } else {
      data.resultTip = ""
    }
    //console.log(result)
  })
} 

</script>

<template>
  <main>
    <table>
      <tbody>
      <tr>
        <!--
        <td width="100">Left</td>
        <td width="100">Right</td>
        -->
        <td colspan="2" style="text-align:left"> 
        <input type="radio" id="folder" value="true" v-model="data.picked" />
        <label for="folder">比对文件夹</label>
        &nbsp; &nbsp; &nbsp;
        <input type="radio" id="files" value="false" v-model="data.picked" />
        <label for="files">比对文件</label>
        </td>
      </tr>
      <tr>
        <!--
        <td></td>
        <td></td>
        -->
        <td width="100%">
          <input id="old" v-model="data.old" autocomplete="off" class="input" type="text"/>
        </td >
        <td width="100">
          <button class="btn" @click="selectOld">源文件</button>
        </td>
      </tr>
      <tr>
      <!--
        <td></td>
        <td></td>
        -->
        <td>
          <input id="new" v-model="data.new" autocomplete="off" class="input" type="text"/>
        </td>
        <td>
          <button class="btn" @click="selectNew">目标文件</button>
        </td>
      </tr>

      <tr>
      <!--
        <td colspan="3" style="text-align:right"> </td>
        -->
        <td colspan="1" style="text-align:right"></td>
        <td>
          <button class="btn" :disabled=data.compareDisabled @click="compare">开始比对</button>
        </td>
      </tr>

      <tr>
        <!--
        <td></td>
        <td></td>
        -->
        <td colspan="1">
          <p> {{ data.resultTip }}</p>
          <textarea v-model="data.result" class="textareaResult"></textarea>
        </td>
        <td></td>
      </tr>
      </tbody>
    </table>

  </main>
</template>

<style scoped>
table {
  margin: 10px;
}
.result {
  height: 20px;
  line-height: 20px;
  margin: 1.5rem auto;
}

.input-box {
  margin: 10px;
}
.btn {
  width: 80px;
  height: 30px;
  line-height: 30px;
  border-radius: 3px;
  border: none;
  margin: 0 0 0 20px;
  padding: 0 8px;
  cursor: pointer;
}

.input-box .btn:hover {
  background-image: linear-gradient(to top, #cfd9df 0%, #e2ebf0 100%);
  color: #333333;
}

.input {
  border: none;
  border-radius: 3px;
  outline: none;
  height: 30px;
  width: 100%;
  padding: 0 6px;
  border: 1px solid #ccc;
  background-color: rgba(240, 240, 240, 1);
  -webkit-font-smoothing: antialiased;
}

.input-box .input:hover {
  /*
  border: none;
  background-color: rgba(255, 255, 255, 1);
  */
}

.input-box .input:focus {
  /*
  border: none;
  background-color: rgba(255, 255, 255, 1);
  */
}

.textareaResult {
  border: 1px solid #ccc;
  border-radius: 3px;
  outline: none;
  line-height: 20px;
  padding: 6px;
  height: 600px;
  width: 100%;
  background-color: rgba(240, 240, 240, 1);
  -webkit-font-smoothing: antialiased;
}
</style>
