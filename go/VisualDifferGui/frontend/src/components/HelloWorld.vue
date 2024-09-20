<script setup>
import {reactive} from 'vue'
import {Greet, SelectOld, SelectOldFolder, SelectNew, MessageBox, CallCompare} from '../../wailsjs/go/main/App'

const data = reactive({
  name: "",
  compareType: false,
  picked: "false",
  compareDisabled: false,
  old: "",
  new: "",
  showResult: false,

  compareObj: {}
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

function test() {
  console.log("test");
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

    console.log(result)
    let json = eval('('+result+')')
    console.log(json)
    console.log(json.Del)

    for(let k in json.Del) {
      console.log(k, json.Del[k]);
    }

    for(let k in json.Change) {
      console.log(k, json.Change[k]);
    }

    data.compareObj = json
    /*
    console.log(json.Sli)
    for(let k in json.Sli) {
      console.log(k, json.Sli.k);
    }
    */
    /*
    console.log(json.CHANGE.length)
    console.log(json.DEL.length)
    console.log(json.NEW.length)
    */
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
          text-sm 
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
        <p 
        v-if = "data.compareObj.Tips != ''"
        class="
            mb-2
            text-indigo-700
            font-semibold
          "> {{ data.compareObj.Tips }}</p>

        <div
            v-if="data.compareObj.Tpo && data.compareObj.Diff"
            class="
            w-full border-2 
            border-indigo-500 
            rounded-md 
            bg-white
            "> 
          <table v-show="data.compareObj.Diff">
              <tbody>
                <tr class="w-full bg-white border-1 border-indigo-500" v-for="item of data.compareObj.Del">
                  <td class="w-1/2 border-1 border-indigo-500  line-through text-rose-700	"> {{ item }} </td>
                  <td class="w-2 border-1 border-indigo-500  font-semibold text-rose-700	"> - </td>
                  <td class="w-1/2 border-1 border-indigo-500"> </td>
                </tr>
                <tr class="w-full bg-white" v-for="item of data.compareObj.Add">
                  <td class="w-1/2"> </td>
                  <td class="w-2 font-semibold text-emerald-900"> + </td>
                  <td class="w-1/2 text-emerald-900"> {{ item }} </td>
                </tr>
                <tr class="w-full bg-white" v-for="val, key of data.compareObj.Change">
                  <td class="w-1/2 text-yellow-700"> {{ key }} </td>
                  <td class="w-2 font-semibold text-yellow-700"> != </td>
                  <td class="w-1/2" @dbclick="test"> {{ val }} </td>
                </tr>
              </tbody>
          </table>
        </div>

        <div v-if="data.compareObj.Tpo==0 && data.compareObj.Diff">
            <div
            v-html="data.compareObj.SingleFileDiff"
            class="
            border-2 
            w-10/12	
            text-left	
            text-wrap	
            border-indigo-500 
            bg-white
            rounded-md 
            p-1.5 
            max-h-96	
            overflow-scroll	
            "></div>
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
