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
  detail: {
    show: true,
    title: "",
    content: "snil<br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br>bdfdsfsafsaf",
  },
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

function detailShow(key) {
  data.detail.show = true
  data.detail.title = key
  data.detail.content = data.compareObj.Change[key]
}

function detailClose() {
  data.detail.show = false
  data.detail.content = ""
}

function compare() {
  if (!(data.old.length && data.new.length)) {
    MessageBox("请提供比对文件或文件夹")
    return
  }

  data.compareDisabled = true
  CallCompare(data.old, data.new).then(result => {
    data.compareDisabled = false

    if (result.length) {
      console.log(result)
      let json = eval('('+result+')')
      console.log(json)
      //console.log(json.Del)

      /*
      for(let k in json.Del) {
        console.log(k, json.Del[k]);
      }

      for(let k in json.Change) {
        console.log(k, json.Change[k]);
      }
      */

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
    }
  })
} 

</script>

<template>
  <main>
    <table class='m-2'>
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
            <input id="old" class="mb-2 w-full h-10 border-2 rounded-md p-1.5 border-indigo-500" v-model="data.old"
              autocomplete="off" type="text" />
          </td>
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
          " @click="selectOld">
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
            <input id="new" class="w-full h-10 border-2 rounded-md p-1.5 border-indigo-500" v-model="data.new"
              autocomplete="off" type="text" />
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
            <button class="
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
          " :disabled=data.compareDisabled @click="compare">开始比对</button>
          </td>
        </tr>

        <tr>
          <td colspan="2">
            <p v-if="data.compareObj.Tips != ''" class="
            mb-2
            text-indigo-700
            font-semibold
          "> {{ data.compareObj.Tips }}</p>

            <div v-if="data.compareObj.Tpo && data.compareObj.Diff" class="
            w-full border-2 
            border-b-none
            border-indigo-500 
            rounded-md 
            bg-white
            ">
              <table v-show="data.compareObj.Diff">
                <thead>
                  <tr class="w-full h-11 bg-indigo-600 border-b-2 border-indigo-500 font-semibold text-white text-lg	">
                    <td class="w-1/2">
                      {{ data.compareObj.Source }}
                    </td>
                    <td class="w-2"> </td>
                    <td class="w-1/2">
                      {{ data.compareObj.Dest }}
                    </td>
                  </tr>
                </thead>
                <tbody>
                  <tr class="w-full h-10 bg-white border-b-2 border-indigo-500 text-rose-700 font-semibold  text-lg"
                    v-for="item of data.compareObj.Del">
                    <td class="line-through"> {{ item }} </td>
                    <td class=""> - </td>
                    <td class=""> </td>
                  </tr>
                  <tr class="w-full h-10 bg-white border-b-2 border-indigo-500 text-emerald-900  font-semibold  text-lg"
                    v-for="item of data.compareObj.Add">
                    <td class=""> </td>
                    <td class=""> + </td>
                    <td class=""> {{ item }} </td>
                  </tr>
                  <tr @dblclick="detailShow(key)" v-for="val, key of data.compareObj.Change" class="cursor-pointer
 w-full h-10 bg-white border-b-2 border-indigo-500 text-yellow-700 font-semibold  text-lg">
                    <td class=""> {{ key }} </td>
                    <td class=""> != </td>
                    <td class=""> {{ key }} </td>
                  </tr>
                </tbody>
              </table>
            </div>

            <div v-if="data.compareObj.Tpo == 0 && data.compareObj.Diff">
              <div v-html="data.compareObj.SingleFileDiff" class="
            border-2 
            w-12/12	
            text-left	
            text-wrap	
            border-indigo-500 
            bg-white
            rounded-md 
            p-1.5 
            min-h-96	
            max-h-96	
            overflow-scroll	
            "></div>
            </div>

          </td>
        </tr>
      </tbody>
    </table>

    <!-- detail view -->
    <div v-if="data.detail.show">
      <div class="
              absolute text-left inset-0 bg-white p-2 border-2
              ">
        <div class="">
          <button class=" 
              min-w-32
              mb-2
              h-11
              border-4 
              rounded-xl
              bg-indigo-500
              sm:w-auto 
              bg-indigo-600 
              hover:bg-indigo-700 
              text-white 
              text-sm 
              font-semibold 
              shadow 
              focus:outline-none 
              cursor-pointer
          " @click="detailClose">返回列表</button>
          <p class="
                  mb-2
                  text-indigo-700
                  font-semibold
                ">{{ data.detail.title }}</p>
        </div>

        <div v-html="data.detail.content" class="
            relative
            border-2 
            w-12/12	
            text-left	
            text-wrap	
            border-indigo-500 
            bg-white
            rounded-md 
            p-1.5 
            h-5/6
            overflow-scroll	
            "></div>
      </div>
    </div>


  </main>
</template>

<style scoped></style>
