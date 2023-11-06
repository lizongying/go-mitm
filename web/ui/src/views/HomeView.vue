<template>
  <div>
    <a-space wrap style="margin-bottom: .5rem">
      <template v-if="isRecord">
        <a-tooltip>
          <template #title>On</template>
          <PoweroffOutlined @click="record" class="hand" style="color: green"/>
        </a-tooltip>
      </template>
      <template v-else>
        <a-tooltip>
          <template #title>Off</template>
          <PoweroffOutlined @click="record" class="hand" style="color: red"/>
        </a-tooltip>
      </template>
      <a-tooltip>
        <template #title>Clear</template>
        <ClearOutlined @click="clear" class="hand"/>
      </a-tooltip>
      <a-tooltip>
        <template #title>Include</template>
        <PlusSquareOutlined @click="filter('Include')" class="hand"/>
      </a-tooltip>
      <a-tooltip>
        <template #title>Exclude</template>
        <MinusSquareOutlined @click="filter('Exclude')" class="hand"/>
      </a-tooltip>
      <a-tooltip>
        <template #title>Proxy</template>
        <GatewayOutlined @click="filter('Proxy')" class="hand"/>
      </a-tooltip>
      <a-tooltip>
        <template #title>Replace</template>
        <DiffOutlined @click="filter('Replace')" class="hand"/>
      </a-tooltip>
      <a-tooltip>
        <template #title>Export</template>
        <ExportOutlined @click="out" class="hand"/>
      </a-tooltip>
      <a-tooltip>
        <template #title>Download Certificate</template>
        <DownloadOutlined @click="certificate" class="hand"/>
      </a-tooltip>
      <a-tooltip>
        <template #title>Question</template>
        <QuestionOutlined @click="filter('Question')" class="hand"/>
      </a-tooltip>
      <span>proxy: {{ lanIp }}:{{ proxyPort }} / {{ internetIp }}:{{ proxyPort }}</span>
    </a-space>
    <a-modal v-model:open="openModal" :title="title" @ok="handleOk">
      <a-textarea v-if="title === 'Include'" v-model:value="include" placeholder="" :rows="4"/>
      <a-textarea v-if="title === 'Exclude'" v-model:value="exclude" placeholder="" :rows="4"/>
      <a-input v-if="title === 'Proxy'" v-model:value="proxy" placeholder=""/>
      <a-space
          v-if="title === 'Replace'"
          v-for="(ruler) in rulers"
          :key="ruler.id"
          style="display: flex; margin-bottom: 8px"
          align="baseline"
      >
        <a-input-group compact>
          <a-input v-model:value="ruler.from" :status="ruler.fromError" :placeholder="ruler.fromPlaceholder"
                   @change="fromChange(ruler)"
                   style="width: 35%">
          </a-input>
          <a-input v-model:value="ruler.to" :status="ruler.toError" :placeholder="ruler.toPlaceholder"
                   @change="toChange(ruler)"
                   style="width: 65%">
            <template #addonBefore>
              <a-select v-model:value="ruler.type" style="width: 93.4px">
                <a-select-option value="https://">https://</a-select-option>
                <a-select-option value="http://">http://</a-select-option>
                <a-select-option value="file://">file://</a-select-option>
              </a-select>
            </template>
            <template #addonAfter>
              <MinusCircleOutlined @click="removeRuler(ruler)"/>
            </template>
          </a-input>
        </a-input-group>
      </a-space>
      <a-form-item v-if="title === 'Replace'">
        <a-button type="dashed" block @click="addRuler">
          <PlusOutlined/>
          Add Ruler
        </a-button>
      </a-form-item>
      <a-list v-if="title === 'Question'" item-layout="horizontal" :data-source="dataQuestion">
        <template #renderItem="{ item }">
          <a-list-item>
            <a-list-item-meta
                :description="item.question"
            >
              <template #title>
                <span>{{ item.ask }}</span>
              </template>
            </a-list-item-meta>
          </a-list-item>
        </template>
      </a-list>
    </a-modal>
    <a-table :columns="columns"
             :data-source="data"
             :pagination="{ hideOnSinglePage: true, pageSize: Infinity }"
             :scroll="{ y: 715 }"
    >
      <template #headerCell="{ column }">
        <template
            v-if="['url', 'method', 'type', 'status', 'time', 'size', 'action'].includes(column.dataIndex)">
        <span style="font-weight: bold">
            {{
            column.dataIndex.charAt(0).toUpperCase()
            + column.dataIndex.slice(1)
          }}
        </span>
        </template>
      </template>
      <template #bodyCell="{ column, text, record}">
        <template v-if="column.dataIndex==='time'">
          {{ text !== undefined ? text / 1000 + 's' : '' }}
        </template>
        <template v-if="column.dataIndex==='size'">
          {{ text !== undefined ? text / 1000 + 'k' : '' }}
        </template>
        <template v-if="['action'].includes(column.dataIndex)">
          <a-space wrap>
            <a-tooltip>
              <template #title>Replay</template>
              <ReloadOutlined @click="replay({
    method: record.method,
    url: record.url,
    req_header: record.req_header,
    req_body: record.req_body,
  })" class="hand"/>
            </a-tooltip>
            <a-tooltip>
              <template #title>More</template>
              <RightOutlined @click="detail(record)" class="hand"/>
            </a-tooltip>
          </a-space>
        </template>
        <template v-if="column.dataIndex === 'status'">
          <a-tag
              :color="text >= 200 && text < 300 ? 'green' : text >= 300 && text < 400  ? 'geekblue' : 'volcano'"
          >
            {{ text }}
          </a-tag>
        </template>
      </template>
      <template
          #customFilterDropdown="{ setSelectedKeys, selectedKeys, confirm, clearFilters, column }"
      >
        <div style="padding: 8px">
          <a-input
              ref="searchInput"
              :placeholder="`Search ${column.dataIndex}`"
              :value="selectedKeys[0]"
              style="width: 188px; margin-bottom: 8px; display: block"
              @change="e => setSelectedKeys(e.target.value ? [e.target.value] : [])"
              @pressEnter="handleSearch(selectedKeys, confirm, column.dataIndex)"
          />
          <a-button
              type="primary"
              size="small"
              style="width: 90px; margin-right: 8px"
              @click="handleSearch(selectedKeys, confirm, column.dataIndex)"
          >
            <template #icon>
              <SearchOutlined/>
            </template>
            Search
          </a-button>
          <a-button size="small" style="width: 90px" @click="handleReset(clearFilters, confirm)">
            Reset
          </a-button>
        </div>
      </template>
      <template #customFilterIcon="{ filtered, column }">
        <template v-if="column.dataIndex === 'url'">
          <search-outlined :style="{ color: filtered ? '#108ee9' : undefined }"/>
        </template>
        <template v-else>
          <filter-outlined :style="{ color: filtered ? '#108ee9' : undefined }"/>
        </template>
      </template>
    </a-table>
    <a-drawer
        v-model:open="open"
        class="custom-class"
        root-class-name="root-class-name"
        :root-style="{ color: 'blue' }"
        style="color: red"
        placement="right"
        size="large"
        :closable="false"
    >
      <a-tabs v-model:activeKey="activeKey">
        <a-tab-pane key="1" tab="Overview">
          <a-collapse v-model:activeKey="activeKeyOverview" :bordered="false">
            <a-collapse-panel key="1" header="Overview">
              <PreviewHeader :header="{'remote addr': message.remote_addr ?? ''}"></PreviewHeader>
            </a-collapse-panel>
            <a-collapse-panel key="2" header="Request Header">
              <PreviewHeader :header="message.req_header"></PreviewHeader>
            </a-collapse-panel>
            <a-collapse-panel key="3" header="Request Cookie">
              <PreviewHeader :header="message.req_cookie"></PreviewHeader>
            </a-collapse-panel>
            <a-collapse-panel key="4" header="Request Body">
              <pre class="language-"><code class="language-">{{ message.req_body ?? '' }}</code></pre>
            </a-collapse-panel>
            <a-collapse-panel key="5" header="Response Header">
              <PreviewHeader :header="message.resp_header"></PreviewHeader>
            </a-collapse-panel>
            <a-collapse-panel key="6" header="Response Cookie">
              <PreviewHeader :header="message.resp_cookie"></PreviewHeader>
            </a-collapse-panel>
            <a-collapse-panel key="7" header="Connection">
              <PreviewHeader :header="message.resp_tls"></PreviewHeader>
            </a-collapse-panel>
          </a-collapse>
        </a-tab-pane>
        <a-tab-pane key="2" tab="Preview">
          <a-collapse v-model:activeKey="activeKeyPreview" :bordered="false">
            <a-collapse-panel key="1" header="raw">
              <pre class="language-"><code class="language-">{{ message.resp_body }}</code></pre>
            </a-collapse-panel>
            <a-collapse-panel key="2" header="hex">
              <pre class="language-"><code class="language-">{{ formatHexDump(message.resp_body) }}</code></pre>
            </a-collapse-panel>
            <a-collapse-panel key="3" header="json">
              <PreviewCode :code="message.resp_body" type="json" line-numbers></PreviewCode>
            </a-collapse-panel>
            <a-collapse-panel key="4" header="html">
              <PreviewCode :code="message.resp_body" type="html" line-numbers></PreviewCode>
            </a-collapse-panel>
          </a-collapse>
        </a-tab-pane>
        <a-tab-pane key="3" tab="cURL">
          <a-collapse v-model:activeKey="activeKeyCurl" :bordered="false">
            <a-collapse-panel key="1" header="curl">
              <PreviewCode
                  :code="toCurl(message.url, message.method, message.req_header, message.req_body)"
                  type="shell" line-numbers></PreviewCode>
            </a-collapse-panel>
          </a-collapse>
        </a-tab-pane>
        <a-tab-pane key="4" tab="javascript">
          <a-collapse v-model:activeKey="activeKeyJs" :bordered="false">
            <a-collapse-panel key="1" header="fetch">
              <PreviewCode
                  :code="toJsFetch(message.url, message.method, message.req_header, message.req_body)"
                  type="javascript" line-numbers></PreviewCode>
            </a-collapse-panel>
            <a-collapse-panel key="2" header="xhr">
              <PreviewCode
                  :code="toJsXhr(message.url, message.method, message.req_header, message.req_body)"
                  type="javascript" line-numbers></PreviewCode>
            </a-collapse-panel>
            <a-collapse-panel key="3" header="ajax">
              <PreviewCode
                  :code="toJsJQuery(message.url, message.method, message.req_header, message.req_body)"
                  type="javascript" line-numbers></PreviewCode>
            </a-collapse-panel>
            <a-collapse-panel key="4" header="axios">
              <PreviewCode
                  :code="toJsAxios(message.url, message.method, message.req_header, message.req_body)"
                  type="javascript" line-numbers></PreviewCode>
            </a-collapse-panel>
          </a-collapse>
        </a-tab-pane>
        <a-tab-pane key="5" tab="python">
          <a-collapse v-model:activeKey="activeKeyPy" :bordered="false">
            <a-collapse-panel key="1" header="request">
              <PreviewCode
                  :code="toPyRequests(message.url, message.method, message.req_header, message.req_body)"
                  type="python" line-numbers></PreviewCode>
            </a-collapse-panel>
            <a-collapse-panel key="2" header="httpx">
              <PreviewCode
                  :code="toPyhttpx(message.url, message.method, message.req_header, message.req_body)"
                  type="python" line-numbers></PreviewCode>
            </a-collapse-panel>
            <a-collapse-panel key="3" header="aiohttp">
              <PreviewCode
                  :code="toPyaiohttp(message.url, message.method, message.req_header, message.req_body)"
                  type="python" line-numbers></PreviewCode>
            </a-collapse-panel>
            <a-collapse-panel key="4" header="urllib">
              <PreviewCode
                  :code="toPyurllib(message.url, message.method, message.req_header, message.req_body)"
                  type="python" line-numbers></PreviewCode>
            </a-collapse-panel>
            <a-collapse-panel key="5" header="urllib3">
              <PreviewCode
                  :code="toPyurllib3(message.url, message.method, message.req_header, message.req_body)"
                  type="python" line-numbers></PreviewCode>
            </a-collapse-panel>
          </a-collapse>
        </a-tab-pane>
        <a-tab-pane key="6" tab="golang">
          <a-collapse v-model:activeKey="activeKeyGo" :bordered="false">
            <a-collapse-panel key="1" header="http.Client">
              <PreviewCode
                  :code="toGoHttpClient(message.url, message.method, message.req_header, message.req_body)"
                  type="go" line-numbers></PreviewCode>
            </a-collapse-panel>
            <a-collapse-panel key="2" header="go-crawler">
              <PreviewCode
                  :code="toGoGoCrawler(message.url, message.method, message.req_header, message.req_body)"
                  type="go" line-numbers></PreviewCode>
            </a-collapse-panel>
          </a-collapse>
        </a-tab-pane>
        <a-tab-pane key="7" tab="java">
          <a-collapse v-model:activeKey="activeKeyJava" :bordered="false">
            <a-collapse-panel key="1" header="HttpClient">
              <PreviewCode :code="toJavaHttpClient(message.url, message.method, message.req_header, message.req_body)"
                           type="java" line-numbers></PreviewCode>
            </a-collapse-panel>
            <a-collapse-panel key="2" header="ApacheHttpClient">
              <PreviewCode
                  :code="toJavaApacheHttpClient(message.url, message.method, message.req_header, message.req_body)"
                  type="java" line-numbers></PreviewCode>
            </a-collapse-panel>
            <a-collapse-panel key="3" header="OKhttp">
              <PreviewCode :code="toJavaOkHttp(message.url, message.method, message.req_header, message.req_body)"
                           type="java" line-numbers></PreviewCode>
            </a-collapse-panel>
            <a-collapse-panel key="4" header="Jsoup">
              <PreviewCode :code="toJavaJsoup(message.url, message.method, message.req_header, message.req_body)"
                           type="java" line-numbers></PreviewCode>
            </a-collapse-panel>
          </a-collapse>
        </a-tab-pane>
        <a-tab-pane key="8" tab="php">
          <a-collapse v-model:activeKey="activeKeyPhp" :bordered="false">
            <a-collapse-panel key="1" header="cURL">
              <PreviewCode
                  :code="toPHPcURL(message.url, message.method, message.req_header, message.req_body)"
                  type="php" line-numbers></PreviewCode>
            </a-collapse-panel>
            <a-collapse-panel key="1" header="Guzzle">
              <PreviewCode
                  :code="toPHPGuzzle(message.url, message.method, message.req_header, message.req_body)"
                  type="php" line-numbers></PreviewCode>
            </a-collapse-panel>
          </a-collapse>
        </a-tab-pane>
      </a-tabs>
    </a-drawer>
  </div>
</template>

<script setup>
import {
  ClearOutlined,
  DiffOutlined,
  DownloadOutlined,
  ExportOutlined,
  FilterOutlined,
  GatewayOutlined,
  MinusCircleOutlined,
  MinusSquareOutlined,
  PlusOutlined,
  PlusSquareOutlined,
  PoweroffOutlined,
  QuestionOutlined,
  ReloadOutlined,
  RightOutlined,
  SearchOutlined
} from '@ant-design/icons-vue';
import {onBeforeMount, reactive, ref} from 'vue'
import {action, event, info} from '@/request/api'
import {formatHexDump} from '@/utils'
import PreviewCode from '../components/PreviewCode.vue'
import PreviewHeader from "../components/PreviewHeader.vue";
import {toCurl} from "@/converter/toCurl";
import {toJsJQuery} from "@/converter/toJsJQuery";
import {toJsXhr} from "@/converter/toJsXhr";
import {toJsFetch} from "@/converter/toJsFetch";
import {toJsAxios} from "@/converter/toJsAxios";
import {toGoHttpClient} from "@/converter/toGoHttpClient";
import {toGoGoCrawler} from "@/converter/toGoGoCrawler";
import {toPHPcURL} from "@/converter/toPHPcURL";
import {toPHPGuzzle} from "@/converter/toPHPGuzzle";
import {toPyRequests} from "@/converter/toPyRequests";
import {toPyurllib3} from "@/converter/toPyurllib3";
import {toPyurllib} from "@/converter/toPyurllib";
import {toPyaiohttp} from "@/converter/toPyaiohttp";
import {toPyhttpx} from "@/converter/toPyhttpx";
import {toJavaHttpClient} from "@/converter/toJavaHttpClient";
import {toJavaJsoup} from "@/converter/toJavaJsoup";
import {toJavaApacheHttpClient} from "@/converter/toJavaApacheHttpClient";
import {toJavaOkHttp} from "@/converter/toJavaOkHttp";
import {notification} from "ant-design-vue";

// record
const isRecord = ref(true)

const record = () => {
  isRecord.value = !isRecord.value
  action({record: isRecord.value}).finally(
      _ => {
        if (isRecord.value) {
          getData()
        }
      }
  )
}

// replay
const replay = record => {
  action({replay: JSON.stringify(record)})
}

// message
const message = ref({})

const open = ref(false);

// detail
const detail = record => {
  message.value = record
  open.value = true;
}

const state = reactive({
  searchText: '',
  searchedColumn: '',
});

const searchInput = ref();
const handleSearch = (selectedKeys, confirm, dataIndex) => {
  confirm();
  state.searchText = selectedKeys[0];
  state.searchedColumn = dataIndex;
};

const handleReset = (clearFilters, confirm) => {
  confirm();
  clearFilters({confirm: true});
  state.searchText = '';
};

const columns = [
  {
    title: 'Url',
    dataIndex: 'url',
    ellipsis: true,
    customFilterDropdown: true,
    onFilter: (value, record) => record.url.toString().toLowerCase().includes(value.toLowerCase()),
    onFilterDropdownOpenChange: visible => {
      if (visible) {
        setTimeout(() => {
          searchInput.value.focus();
        }, 100);
      }
    },
  },
  {
    title: 'Type',
    dataIndex: 'type',
    ellipsis: true,
    width: 200,
    filters: [
      {
        text: 'plain',
        value: 'plain',
      },
      {
        text: 'html',
        value: 'html',
      },
      {
        text: 'json',
        value: 'json',
      },
      {
        text: 'javascript',
        value: 'javascript',
      },
      {
        text: 'css',
        value: 'css',
      },
      {
        text: 'image',
        value: 'image',
      },
      {
        text: 'video',
        value: 'video',
      },
      {
        text: 'font',
        value: 'font',
      },
    ],
    onFilter: (value, record) => {
      if (record.type === undefined) {
        return false
      }
      return record.type.indexOf(value) > -1
    },
  },
  {
    name: 'Method',
    dataIndex: 'method',
    width: 100,
    filters: [
      {
        text: 'GET',
        value: 'GET',
      },
      {
        text: 'POST',
        value: 'POST',
      },
      {
        text: 'HEAD',
        value: 'HEAD',
      },
      {
        text: 'PUT',
        value: 'PUT',
      },
      {
        text: 'DELETE',
        value: 'DELETE',
      },
      {
        text: 'PATCH',
        value: 'PATCH',
      },
      {
        text: 'OPTIONS',
        value: 'OPTIONS',
      },
      {
        text: 'CONNECT',
        value: 'CONNECT',
      },
      {
        text: 'TRACE',
        value: 'TRACE',
      }
    ],
    onFilter: (value, record) => record.method === value,
  },
  {
    title: 'Status',
    dataIndex: 'status',
    width: 100,
    filters: [
      {
        text: '200',
        value: 200,
      },
      {
        text: '300',
        value: 300,
      },
      {
        text: '400',
        value: 400,
      },
      {
        text: '500',
        value: 500,
      },
    ],
    onFilter: (value, record) => record.status >= value && record.status < value + 100,
  },
  {
    title: 'Time',
    dataIndex: 'time',
    width: 100,
    sorter: (a, b) => a.time - b.time,
    sortDirections: ['descend', 'ascend'],
  },
  {
    title: 'Size',
    dataIndex: 'size',
    width: 100,
    sorter: (a, b) => a.size - b.size,
    sortDirections: ['descend', 'ascend'],
  },
  {
    title: 'Action',
    dataIndex: 'action',
    width: 100,
  },
];
const data = ref([]);

const maxRow = 50

const getData = () => {
  const es = event()
  es.onopen = _ => {
    // console.log('es open');
  };
  es.onmessage = event => {
    // console.log(event.data);
    data.value.unshift(JSON.parse(event.data))
    if (data.value.length > maxRow) {
      data.value = data.value.slice(0, maxRow)
    }
  };
  es.onerror = event => {
    // console.log('es error', event)
    if (event.readyState === EventSource.CLOSED) {
      // console.log('event was closed')
    }
  };
  es.addEventListener('close', _ => {
    // console.log('es close')
    es.close()
    isRecord.value = false
  })
}

const activeKey = ref('1');
const activeKeyOverview = ref([2])
const activeKeyPreview = ref([1])
const activeKeyCurl = ref([1])
const activeKeyJs = ref([4])
const activeKeyPy = ref([1])
const activeKeyGo = ref([1])
const activeKeyPhp = ref([1])
const activeKeyJava = ref([1])

// clear data
const clear = () => {
  data.value = []
}

const certificate = () => {
  const file = new File([`-----BEGIN CERTIFICATE-----
MIIC4TCCAcmgAwIBAgIBATANBgkqhkiG9w0BAQsFADASMRAwDgYDVQQDEwdHTy1N
SVRNMB4XDTIzMDkxMTE0MzQ0NVoXDTMzMDkxMTE0MzQ0NVowEjEQMA4GA1UEAxMH
R08tTUlUTTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAK0wrAupDXDu
hQykEhV6WTNWL5GW9sA6k0dGM5lPceD6zvt3eBS6rfLboq4Bw4RQXXv+BQXVDiSa
v/BH62RCxp/3EkSsj5JbeksZdxO42JqJ1hzVzEGwmNDBYy2vjKEnyS7z/PK8LOwS
jDiaWkb+dc/7Ixomb50e3eT5WCxRoL/HVWaePSjV8FMmt2anDiicJEo0VJoF4A18
Qme3ueXvFL6tbpCIHMr/BwpV7suuRS5/KyLlI12TOYQoomJn86Ru+TOENrD/Gb4B
u9nXqAq3gr5Ke+akKxTNz/8Mxz8+08xGDQE/4YhkFy28v0kgwFmZOVLdzVNZSoRE
lWPswBIpD6cCAwEAAaNCMEAwDgYDVR0PAQH/BAQDAgKkMA8GA1UdEwEB/wQFMAMB
Af8wHQYDVR0OBBYEFMDqagzu5CwpAlDc8RXIX5Dl8pZIMA0GCSqGSIb3DQEBCwUA
A4IBAQCmbwUN4P5aqoGL0HwI9ZhX3PuEtixNL1+Ev33JvhbezqtvHFBjrZK/7TkF
wRX9nHtcKP3vKT8Y+n/P5x9KXvicYwrLm2wF6iXhhdtUfFkbZi5ZtIurSWOaJbzV
Sx2AOl6t2yrcR2wljme0dqVuPRLqMehglUhC/yQYqkg5CDnR6VkMCKcMb7UniVfX
hVECMhnS/rMfXdNmnuHJJBUeD+FW67kfng4GyHPDWokziDtedaP/YDLL3k9g8ctG
HXZ5VDduL8tIPTDbsRZIoXKfUQXi0Otr7XfruOTESuiy+ux8TQeAT2piLtchcRO0
4FnakJzTJLa0wpgGR0iDPRNUPSML
-----END CERTIFICATE-----`], 'ca.crt', {
    type: 'text/csv',
  });

  const tmpLink = document.createElement('a')
  const objectUrl = URL.createObjectURL(file)

  tmpLink.href = objectUrl
  tmpLink.download = file.name
  tmpLink.click()
  URL.revokeObjectURL(objectUrl)
}

// export data
const out = () => {
  const text = ['Url,Type,Method,Status,Time,Size']
  data.value.forEach((i) => {
    text.push([i.url, i.type, i.method, i.status, i.time, t.size].join(','))
  })
  const file = new File([text.join('\n')], 'go_mitm.csv', {
    type: 'text/csv',
  });

  const tmpLink = document.createElement('a')
  const objectUrl = URL.createObjectURL(file)

  tmpLink.href = objectUrl
  tmpLink.download = file.name
  tmpLink.click()
  URL.revokeObjectURL(objectUrl)
}

const title = ref('')
const include = ref('')
const exclude = ref('')
const proxy = ref('')

const rulers = reactive([])
const removeRuler = item => {
  const index = rulers.indexOf(item);
  if (index !== -1) {
    rulers.splice(index, 1);
  }
  doReplace()
};
const addRuler = () => {
  rulers.push({
    from: '',
    fromError: '',
    fromPlaceholder: '',
    to: '',
    toError: '',
    toPlaceholder: '',
    type: 'https://',
    id: Date.now(),
  });
};
const fromChange = ruler => {
  if (ruler.from !== '') {
    ruler.fromError = ''
  } else {
    ruler.fromError = 'error'
  }
}
const toChange = ruler => {
  if (ruler.to !== '') {
    ruler.toError = ''
  } else {
    ruler.toError = 'error'
  }
}

const openModal = ref(false);
const filter = (t) => {
  title.value = t
  openModal.value = true;
}
const handleOk = () => {
  if (title.value === 'Include') {
    let inStr = include.value.split('\n').map(v => {
      return v.trim()
    }).filter(v => {
      return v !== ''
    }).join(';')
    if (inStr === "") {
      inStr = "-"
    }
    action({
      include: inStr
    }).finally(openModal.value = false)
  }
  if (title.value === 'Exclude') {
    let exStr = exclude.value.split('\n').map(v => {
      return v.trim()
    }).filter(v => {
      return v !== ''
    }).join(';')
    if (exStr === "") {
      exStr = "-"
    }
    action({
      exclude: exStr
    }).finally(openModal.value = false)
  }
  if (title.value === 'Proxy') {
    let proxyStr = proxy.value.trim()
    if (proxyStr === "") {
      proxyStr = "-"
    }
    action({
      proxy: proxyStr
    }).finally(openModal.value = false)
  }
  if (title.value === 'Replace') {
    doReplace()
  }
}

const doReplace = () => {
  let arr = []
  let err = false
  rulers.forEach((v) => {
    if (v.from === "") {
      v.fromError = 'error'
      err = true
    }
    if (v.to === "") {
      v.toError = 'error'
      err = true
    }
    if (err) {
      return
    }
    arr.push(`${v.from},${v.type}${v.to}`)
  })
  if (err) {
    notification.warning({
      message: 'warning'
    })
    return
  }
  let show = arr.join(";").trim()
  if (show === "") {
    show = "-"
  }
  action({
    replace: show
  }).finally(() => {
        notification.success({
          message: 'success'
        })
      }
  )
}

const lanIp = ref('')
const internetIp = ref('')
const proxyPort = ref(0)
onBeforeMount(() => {
  info().then(response => {
    const info = response.data
    // console.log('info', info)
    isRecord.value = info.record
    include.value = info.include ? info.include.join('\n') : ''
    exclude.value = info.exclude ? info.exclude.join('\n') : ''
    lanIp.value = info.lan_ip ? info.lan_ip : ''
    proxyPort.value = info.proxy_port ? info.proxy_port : ''
    q4.value = `export HTTP_PROXY=http://${lanIp.value}:${proxyPort.value} && export HTTPS_PROXY=http://${lanIp.value}:${proxyPort.value}`
    q3.value = `curl "https://httpbin.org/get" -x http://${lanIp.value}:${proxyPort.value} --cacert ./ca.crt`
    internetIp.value = info.internet_ip ? info.internet_ip : ''
  }).finally(_ => {
    if (isRecord.value) {
      getData()
    }
  })
})

const q3 = ref('')
const q4 = ref('')
const dataQuestion = reactive([
  {
    ask: '为什么请求没有被记录？',
    question: '1. 请检查Exclude中是否被排除，默认排除目标域名是localhost、127.0.0.1的请求。'
  },
  {
    ask: '支持https请求吗？',
    question: '请在Download Certificate中下载证书，安装并信任。'
  },
  {
    ask: '如何在curl中使用代理？',
    question: q3
  },
  {
    ask: '如何在terminal中使用代理？',
    question: q4
  },
]);

</script>

<style scoped>
.hand {
  cursor: pointer
}
</style>