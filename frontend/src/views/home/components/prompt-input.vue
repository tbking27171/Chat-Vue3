<template>
  <a-space class="prompt-container" direction="horizontal" style="width: 100%;">
    <div class="prompt-input-div">
      <div class="prompt-extension">
        <a-select v-model="curSysPrompt.id" :style="{width:'160px'}" placeholder="Please select ..."
                  @change="handleSysPromptChange">
          <a-option v-for="(item, index) in syProList" :key="index" :value="item.id" :label="item.label"/>
        </a-select>
        <a-apace :style="{marginLeft: '10px'}">
          <span v-if="curSysPrompt.is_sys === 1">{{ $t('common.system') }}:</span>
          <span v-if="curSysPrompt.type === 2 && curSysPrompt.id !== 2">{{ curSysPrompt.desc }}</span>
          <a-input v-model="curSysPrompt.desc" :style="{width:'400px', marginLeft: '10px'}"
                   v-if="curSysPrompt.type === 1"
                   :placeholder="$t('common.prompt.input.system')"/>

          <a-select v-if="curSysPrompt.type === 2 && curSysPrompt.id === 3" v-model="curSysPrompt.language"
                    :style="{width:'120px', marginLeft: '10px'}"
                    placeholder="Please select ...">
            <a-option v-for="(item, index) in curSysPrompt.extra" :key="index" :value="item.label" :label="item.label"/>
          </a-select>
          <div v-if="curSysPrompt.id === 1 || curSysPrompt.id === 2"
               style="position: absolute; right: 10px; top: 10px; width: 230px;">
            <a-select v-model="currReply" :style="{width: '100%'}">
              <a-option v-for="(item, index) in replyPreList" :key="index" :value="item" :label="item"/>
            </a-select>
          </div>
        </a-apace>
      </div>
      <a-textarea v-model="promptValue" class="prompt-textarea" :placeholder="$t('common.prompt.input.placeholder')"
                  @keydown="handleKeyDownSend"/>
      <a-button class="prompt-btn" size="large" @click="handleSend"
                :loading="sendLoading">
        <template #icon>
          <icon-send size="18"/>
        </template>
      </a-button>
    </div>

    <!-- stop chat stream -->
    <div class="backoff-container" v-if="checkOffFlag">
      <a-button size="medium" @click="handleBreakOffChat">
        <template #icon>
          <icon-record-stop/>
        </template>
        {{ $t('common.breakoff') }}
      </a-button>
    </div>
  </a-space>
</template>
<script setup>
import {reactive, ref, watch} from "vue";
import {BreakOffChatStream} from "../../../../wailsjs/go/chat/Service.js";
import {Notification} from "@arco-design/web-vue";
import {useI18n} from "vue-i18n";
import {LOCALE_OPTIONS} from "@/locale";

const {t} = useI18n();

const props = defineProps({
  value: {
    type: String,
    default: ""
  },
  loading: {
    type: Boolean,
    default: false
  },
  checkoff: {
    type: Boolean,
    default: false
  }
})
const syProList = reactive([
  {
    id: 1,
    label: t('common.prompt.select.title1'),
    desc: t('common.prompt.select.title1.desc'),
    type: 2, // 1: chat  2：answer
    is_sys: 2, // 1: is system
    extra: []
  }, {
    id: 2,
    label: t('common.prompt.select.title2'),
    desc: t('common.prompt.select.title2.desc'),
    type: 1,
    is_sys: 1,
    extra: []
  }, {
    id: 3,
    label: t('common.prompt.select.title3'),
    desc: t('common.prompt.select.title3.desc'),
    type: 2,
    is_sys: 2,
    extra: LOCALE_OPTIONS
  }, {
    id: 4,
    label: t('common.prompt.select.title4'),
    desc: t('common.prompt.select.title4.desc'),
    type: 2,
    is_sys: 2,
    extra: []
  }, {
    id: 5,
    label: t('common.prompt.select.title5'),
    desc: t('common.prompt.select.title5.desc'),
    type: 2,
    is_sys: 2,
    extra: []
  }, {
    id: 6,
    label: t('common.prompt.select.title6'),
    desc: t('common.prompt.select.title6.desc'),
    type: 2,
    is_sys: 2,
    extra: []
  }, {
    id: 7,
    label: t('common.prompt.select.title7'),
    desc: t('common.prompt.select.title7.desc'),
    type: 2,
    is_sys: 2,
    extra: []
  }, {
    id: 8,
    label: t('common.prompt.select.title8'),
    desc: t('common.prompt.select.title8.desc'),
    type: 2,
    is_sys: 2,
    extra: []
  },
]);
const promptValue = ref(props.value);
let curSysPrompt = reactive({
  id: 1,
  label: t('common.prompt.select.title1'),
  desc: t('common.prompt.select.title1.desc'),
  type: 2,
  is_sys: 2,
  extra: [],
  language: '',
});
const sendLoading = ref(props.loading);
const checkOffFlag = ref(props.checkoff);
const replyPreList = reactive([t('common.prompt.input.replynormal'), t('common.prompt.input.replymarkdown'), t('common.prompt.input.replyform')]);
const currReply = ref(t('common.prompt.input.replynormal'));

watch(() => props.value, () => {
  promptValue.value = props.value;
})
watch(() => props.loading, () => {
  sendLoading.value = props.loading;
})
watch(() => props.checkoff, () => {
  checkOffFlag.value = props.checkoff;
})
const emits = defineEmits(['add:chat'])

const handleSend = () => {
  if (promptValue.value == "" || sendLoading.value == true) {
    return;
  }
  const sendPromptList = assemblePrompt();
  sendLoading.value = true;
  emits('add:chat', sendPromptList, curSysPrompt, sendLoading.value)
  promptValue.value = "";
}

const assemblePrompt = () => {
  let promptList = [];
  let currReplyContent = "," + currReply.value
  if (currReply.value === t('common.prompt.input.replynormal')) {
    currReplyContent = "";
  }
  switch (curSysPrompt.type) {
    case 1:
      promptList = [{
        role: 'system',
        prefix: '',
        content: curSysPrompt.desc.trim(),
        reply_content: currReplyContent
      }, {
        role: 'user',
        prefix: '',
        content: promptValue.value.trim(),
        reply_content: currReplyContent
      }];
      break;
    case 2:
      let systemPromptValue = curSysPrompt.desc;
      if (curSysPrompt.type === 2 && curSysPrompt.id === 3) {
        systemPromptValue += curSysPrompt.language + ":";
      }
      if (curSysPrompt.type === 2 && curSysPrompt.id === 1) {
        systemPromptValue = "";
      }
      promptList = [{
        role: 'user',
        prefix: systemPromptValue,
        content: promptValue.value.trim(),
        reply_content: currReplyContent
      }]
      break;
  }
  return promptList
}
const handleKeyDownSend = (event) => {
  if (event.ctrlKey && event.keyCode === 13) {
    handleSend()
  }
}
const handleBreakOffChat = () => {
  BreakOffChatStream()
}
const handleSysPromptChange = (e) => {
  if (e === 2) {
    Notification.info({
      content: t('common.prompt.input.chatTips'),
      duration: 3000
    })
  }
  const promptInfo = syProList.filter((item) => {
    return e === item.id
  });
  if (promptInfo.length > 0) {
    curSysPrompt.id = promptInfo[0].id;
    curSysPrompt.label = promptInfo[0].label;
    curSysPrompt.desc = promptInfo[0].desc;
    curSysPrompt.type = promptInfo[0].type;
    curSysPrompt.is_sys = promptInfo[0].is_sys;
    curSysPrompt.extra = promptInfo[0].extra;
  }
  if (curSysPrompt.type === 2 && curSysPrompt.id === 3) {
    curSysPrompt.language = 'English';
  }
}
</script>

<style scoped>
.prompt-container {
  position: relative;
}

.prompt-container .prompt-input-div {
  height: 150px;
  width: 100%;
  padding: 10px;
  background: var(--color-bg-1);
  position: relative;
  border-top: 1px solid var(--color-neutral-2);
}

.prompt-container .prompt-input-div .prompt-textarea {
  height: 110px;
  margin-top: 10px;
  border-radius: 6px;
}

.prompt-container .prompt-extension {
}

.prompt-container .prompt-btn {
  border-radius: 6px;
  width: 45px;
  position: absolute;
  right: 10px;
  bottom: 10px;
  background: rgb(var(--purple-5));
  color: #fff;
}

.backoff-container {
  position: absolute;
  bottom: 175px;
  left: 10px;
}

.prompt-container :deep(.arco-textarea-wrapper) {
  background: none;
  border: none;
}
</style>