<script setup>
import {reactive} from 'vue'

import {Greet, SelectOld, SelectOldFolder, SelectNew, MessageBox, CallCompare} from '../../wailsjs/go/main/App'


const data = reactive({
  lang: "plaintext",
  prev: "",
  current: "",
  name: "",
  compareType: false,
  picked: "false",
  compareDisabled: false,
  compareBtnText: "开始比对 ",
  old: "",
  new: "",
  detail: {
    show: false,
    title: "",
    content: "",
    //content: "snil<br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br><br>bdfdsfsafsaf",
  },
  compareObj: {}
})


data.prev = `var a1 = {
  "name": "vue-diff",
  "version": "0.0.0",
  "description": "Vue diff viewer",
  "private": true,
  "peerDependencies": {
    "vue": "^3.0.0"
  }
}`;


data.current = `const b2 = {
  "name": "vue-diff",
  "version": "1.0.0",
  "description": "Vue diff viewer",
  "scripts": {
    "test": "vue-cli-service test:unit"
    "lint": "vue-cli-service lint"
  },
  "peerDependencies": {
    "vue": "^3.0.0"
  }
}`;
function selectType(type) {
  if (type == 0) {
    data.compareType = false
  } else {
    data.compareType = true
  }
}

function selectOld() {
  /*
  if (data.picked == "false") {
    data.compareType = false
  } else {
    data.compareType = true
  }
  */
  SelectOld(data.compareType).then(result => {
    if (result.length) {
      data.old = result
    }
  })
}

function selectNew() {
  /*
  if (data.picked == "false") {
    data.compareType = false
  } else {
    data.compareType = true
  }
  */

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
  data.compareBtnText =  "处理中...",
  CallCompare(data.old, data.new).then(result => {
    data.compareBtnText =  "开始比对",
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
      data.prev = json.Old
      data.current = json.New
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
    <div class='px-2 pb-0 bg-white h-full h-dvh h-screen'>
      <table class='w-full'>
        <tbody>
          <tr>
            <td colspan="2" style="text-align:left" class="align-middle">
              <button v-if="data.compareType" class="
          min-w-32
          h-10
          border-1 
          rounded-md
          bg-gray-700
          text-white
          font-semibold
          w-full
          sm:w-auto 
          hover:bg-gray-700 
          hover:text-white
          shadow 
          focus:outline-none 
          cursor-pointer
          mr-5 
              ">比对文件夹</button>

              <button v-if="data.compareType == 0" class=" 
          min-w-32
          h-10
          border-1 
          rounded-md
          text-gray-700
          bg-gray-200 
          w-full
          sm:w-auto 
          hover:bg-gray-700 
          hover:text-white
          font-semibold 
          shadow 
          focus:outline-none 
          cursor-pointer
          mr-5
          " @click="selectType(1)">比对文件夹</button>

              <button v-if="data.compareType" class=" 
          min-w-32
          mt-2
          mb-2
          h-10
          rounded-md
          w-full
          sm:w-auto 
          bg-gray-200 
          text-gray-700
          hover:bg-gray-700 
          hover:text-white
          font-semibold 
          shadow 
          focus:outline-none 
          cursor-pointer
          " @click="selectType(0)">比对文件</button>

              <button v-if="!data.compareType" class=" 
          min-w-32
          mt-2
          mb-2
          h-10
          rounded-md
          w-full
          sm:w-auto 
          text-white
          bg-gray-700 
          hover:bg-gray-700 
          hover:text-white
          font-semibold 
          shadow 
          focus:outline-none 
          cursor-pointer
          " @click="selectType(0)">比对文件</button>
              <!--
              <div class="align-top md:align-top ">
                <input type="radio" id="folder" value="true" v-model="data.picked" />
                <label for="folder">&nbsp;比对文件夹</label>

                &nbsp; &nbsp; &nbsp;

                <input type="radio" id="files" value="false" v-model="data.picked" />
                <label for="files"> &nbsp;比对文件</label>
              </div>
            -->
            </td>
          </tr>
          <tr>
            <td class="w-11/12">
              <input id="old"
                class="mb-2 w-full h-10 border-2 hover:border-gray-700 focus:border-gray-700 rounded-md p-1.5 border-gray-300"
                v-model="data.old" autocomplete="off" type="text" />
            </td>
            <td class="w-1/12 text-right">
              <button class="
          min-w-32
          ml-4
          mb-2
          h-10
          border-1 
          rounded-md
          text-gray-700
          font-semibold
          w-full
          sm:w-auto 
          bg-gray-200 
          hover:bg-gray-700 
          hover:text-white
          font-semibold 
          shadow 
          focus:outline-none 
          cursor-pointer
          " @click="selectOld">
                源文件
              </button>
            </td>
          </tr>
          <tr>
            <td>
              <input id="new"
                class="w-full h-10 border-2 hover:border-gray-700 focus:border-gray-700 rounded-md p-1.5 border-gray-300"
                v-model="data.new" autocomplete="off" type="text" />
            </td>
            <td class="text-right">
              <button class=" 
          min-w-32
          mt-2
          mb-2
          h-10
          border-1 
          rounded-md
          text-gray-700
          font-semibold
          w-full
          sm:w-auto 
          bg-gray-200 
          hover:bg-gray-700 
          hover:text-white
          font-semibold 
          shadow 
          focus:outline-none 
          cursor-pointer
          " @click="selectNew">目标文件</button>
            </td>
          </tr>

          <tr>
            <td colspan="1" class="text-right"></td>
            <td class="text-right">
              <button class="
            min-w-32
          mt-2
          mb-2
          h-10
          border-1 
          rounded-md
          text-gray-700
          font-semibold
          w-full
          sm:w-auto 
          bg-gray-200 
          enabled:hover:bg-gray-700 
          enabled:hover:text-white
          font-semibold 
          shadow 
          focus:outline-none 
          cursor-pointer
          disabled:opacity-75
          " :disabled=data.compareDisabled @click="compare">{{ data.compareBtnText }}</button>
            </td>
          </tr>

          <tr>
          <td colspan=2>
          <!-- TD_DIFF -->
<div class="w-full border-2 border-gray-300 rounded-md" >
        <Diff
            mode="split"
            theme="light"
            :language="data.lang"
            :prev="data.prev"
            :current="data.current"
            input-delay="0"
            virtual-scroll="{ height: 500, lineMinHeight: 24, delay: 100 }"
          />
</div>
          </td>
          </tr>

          <tr>
            <td colspan="2">
            <!-- TD_TIP -->
              <p v-if="data.compareObj.Tips != ''" class="
            mb-2
            text-indigo-700
            font-semibold
            h-full
          "> {{ data.compareObj.Tips }}</p>

              <div v-if="data.compareObj.Tpo && data.compareObj.Diff" class="
              border border-b-none border-gray-300
              w-full
              ">
                <table class="w-full border-collapse" v-show="data.compareObj.Diff">
                  <thead>
                    <tr class="h-11 bg-gray-200 font-semibold text-gray-900 text-lg	">
                      <td width="200" class="max-w-60 truncate">
                        {{ data.compareObj.Source }}
                      </td>
                      <td class="w-2"> </td>
                      <td width="200" class="max-w-60 truncate">
                        {{ data.compareObj.Dest }}
                      </td>
                    </tr>
                  </thead>
                </table>
              </div>

              <!-- TD_DIR -->
              <div v-if="data.compareObj.Tpo && data.compareObj.Diff" class="
              border border-b-none border-t-none border-gray-300
              w-full
              h-svh
              max-h-96
              overflow-auto
              ">
                <table class="w-full border-collapse" v-show="data.compareObj.Diff">
                  <tbody>
                    <tr class="h-10 border border-t-none border-gray-200 text-rose-700 font-semibold  text-lg"
                      v-for="item of data.compareObj.Del">
                      <td class="max-w-60 max-w-11/12 line-through truncate" :title="item"> {{ item }} </td>
                      <td class=""> - </td>
                      <td class=""> </td>
                    </tr>
                    <tr class="h-10  border border-gray-200 text-emerald-900  font-semibold  text-lg"
                      v-for="item of data.compareObj.Add">
                      <td class=""> </td>
                      <td class=""> + </td>
                      <td class="max-w-60 truncate" :title="item"> {{ item }} </td>
                    </tr>
                    <tr class="cursor-pointer h-10  border border-gray-200 text-yellow-700 font-semibold  text-lg"
                      @dblclick="detailShow(key)" v-for="val, key of data.compareObj.Change">
                      <td class="max-w-60 truncate" :title="key"> {{ key }} </td>
                      <td class=""> != </td>
                      <td class="max-w-60 truncate" :title="key"> {{ key }} </td>
                    </tr>
                  </tbody>
                </table>
              </div>

              <!-- TD_FILE -->
              <div v-if="data.compareObj.Tpo == 0 && data.compareObj.Diff">
                <div v-html="data.compareObj.SingleFileDiff" class="
            border-2
            w-12/12	
            text-left	
            text-wrap	
            border-gray-500 
            bg-white
            rounded-md 
            p-1.5 
            overflow-scroll	
            max-h-96
            h-screen
            "></div>
              </div>

            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- detail view -->
    <!-- TD_DIR_FILE -->
    <div v-if="data.detail.show">
      <div class="z-50 absolute top-0 left-0 h-20 text-left inset-0 bg-white p-2 border border-10 
      border-b-gray-300
      ">
        <div class="">
          <button class=" 
              h-10
              min-w-32
              rounded-md
              sm:w-auto 
              bg-gray-200 
              text-black 
              hover:bg-gray-700 
              hover:text-white
              text-sm 
              font-semibold 
              cursor-pointer 
          " @click="detailClose">返回列表</button>
          <p class="h-10 mt-1 align-middle text-indigo-700 font-semibold ">{{ data.detail.title }}</p>
        </div>
      </div>
    </div>

    <div v-if="data.detail.show">
      <div class="z-40 absolute bg-white top-20 left-0 w-full h-full">
        <div v-html="data.detail.content" class="
            space-x-0.5
            m-2
            p-2
            border-2 
            text-left	
            text-wrap	
            border-gray-500 
            bg-white
            rounded-md 
            w-11/12
            md:w-auto
            h-5/6
            min-h-content
            overflow-scroll	
            "></div>
      </div>
    </div>


  </main>
</template>

<!--<style scoped></style>-->
<style lang="scss">
.vue-diff-theme-custom {
  @import 'highlight.js/scss/vs2015.scss'; // import theme

  background-color: #000; // Set background color
}
</style>
