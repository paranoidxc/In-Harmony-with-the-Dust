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
  showResult: false,
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
    data.showResult = false
    data.resultTip = ""

    if (result != "-") {
      data.result = result
      data.showResult = true
    }
    if (result.length == 0) {
      data.resultTip = "文件内容相同"
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
        <td colspan="2" style="text-align:left" class="align-middle"> 
        <div class="align-top md:align-top ">
          <input type="radio" id="folder" value="true" v-model="data.picked" />
          <label for="folder">&nbsp;比对文件夹</label>

          &nbsp; &nbsp; &nbsp;

          <input type="radio" id="files" value="false" v-model="data.picked" />
          <label for="files"> &nbsp;比对文件</label>
        </div>
        </td>
      </tr>
      <tr>
        <!--
        <td></td>
        <td></td>
        -->
        <td width="90%">
          <input id="old" class="mb-2 w-full h-10 border-2 rounded-md p-1.5 border-indigo-500" 
          v-model="data.old" autocomplete="off" type="text"/>
        </td >
        <td width="100px">
          <button class="
          min-w-32
          mb-2 h-11
          border-4 
          rounded-xl
          bg-indigo-500
          text-white
          font-semibold
          hover:bg-indigo-700

          w-full sm:w-auto 
          bg-indigo-600 
          hover:bg-indigo-700 
          disabled:bg-indigo-300 
          dark:disabled:bg-indigo-800 
          text-white 
          dark:disabled:text-indigo-400 
          text-sm font-semibold 
          rounded-md 
          shadow 
          focus:outline-none 
          cursor-pointer
          "
          @click="selectOld">
          源文件
          </button>
        </td>
      </tr>
      <tr>
      <!--
        <td></td>
        <td></td>
        -->
        <td>
          <input id="new" 
          class="w-full h-10 border-2 rounded-md p-1.5 border-indigo-500"
          v-model="data.new" autocomplete="off" type="text"/>
        </td>
        <td>
          <button class=" 
          mt-2

          min-w-32
          mb-2 h-11
          border-4 
          rounded-xl
          bg-indigo-500
          text-white
          font-semibold
          hover:bg-indigo-700

          w-full sm:w-auto 
          bg-indigo-600 
          hover:bg-indigo-700 
          disabled:bg-indigo-300 
          dark:disabled:bg-indigo-800 
          text-white 
          dark:disabled:text-indigo-400 
          text-sm font-semibold 
          rounded-md 
          shadow 
          focus:outline-none 
          cursor-pointer
          " @click="selectNew">目标文件</button>
        </td>
      </tr>

      <tr>
      <!--
        <td colspan="3" style="text-align:right"> </td>
        -->
        <td colspan="1" style="text-align:right"></td>
        <td>
          <button 
          class="
          min-w-32
          mb-2 h-11
          border-4 
          rounded-xl
          bg-indigo-600 
          text-white
          font-semibold
          hover:bg-indigo-700
          shadow 
          "
           :disabled=data.compareDisabled @click="compare">开始比对</button>
        </td>
      </tr>

      <tr>
        <!--
        <td></td>
        <td></td>
        -->
        <td colspan="2">
          <p class="
          mb-2
          text-indigo-500
          font-semibold
          "> {{ data.resultTip }}</p>
          <div v-if="data.showResult">
          <textarea v-model="data.result" 
          class="w-full h-dvh 
          w-full border-2 rounded-md p-1.5 border-indigo-500 
          "></textarea>
          </div>
        </td>
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
