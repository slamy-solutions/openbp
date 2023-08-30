<template>
    <q-page class="q-pa-md row">
      <div class="col-2 q-pr-md">
        <MenuComponent selected="oauth"/>
      </div>
      <div class="col-10 row">
        <q-list class="full-width">
          <div v-for="provider of providers" v-bind:key="provider.label" class="row">
          <div :class="selectedProviderLabel == provider.label ? 'col bg-secondary q-mt-xs q-mb-xs': 'col q-mt-xs q-mb-xs'" style="max-width: 8px; border-top: solid 1px; border-bottom: solid 1px; border-left: solid 1px;"></div>
          <q-item
            :class="selectedProviderLabel == provider.label ? 'col bg-grey-2' : 'col'"
            :clickable="selectedProviderLabel != provider.label"
            :v-ripple="selectedProviderLabel != provider.label"
            style="width: 400px; border-top: solid 1px; border-bottom: solid 1px; border-right: solid 1px;"
            class="q-mt-xs q-mb-xs"
            @click="selectProvider(provider)"
          >
            <q-item-section avatar>
              <q-icon :name="provider.icon" />
            </q-item-section>
            <q-item-section>{{ provider.label }}</q-item-section>
            <q-item-section side>
              <q-icon v-if="provider.enabled" color="positive" name="check_circle" />
              <q-icon v-else name="cancel" />
            </q-item-section>
          </q-item>
          </div>
        </q-list>
        
        <q-card
          flat
          class="bg-transparent q-mt-lg q-pt-lg full-width "
          style="border-top: #ccc solid 1px;"
        >
          <q-item tag="label" v-ripple style="border: #ccc solid 1px;">
        <q-item-section>
          <q-item-label>Provider enabled</q-item-label>
          <q-item-label caption>Allow notification</q-item-label>
        </q-item-section>
        <q-item-section side top>
          <q-toggle color="positive" v-model="selectedProviderEnabled" val="friend" />
        </q-item-section>
      </q-item>

          <q-input outlined square label="Client ID" v-model="selectedProviderClientID" class="q-mt-sm"></q-input>
          <q-input outlined square label="Client secret" v-model="selectedProviderSecret" type="password" class="q-mt-sm"></q-input>
          
          <q-btn label="update configuration" outline class="full-width q-mt-md" @click="updateProviderConfig"></q-btn>
        </q-card>

      </div>
    </q-page>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { useRoute } from 'vue-router';
import api from 'src/boot/api';
import { OAuthProviderConfig, OAuthProviderName } from 'src/boot/api/accessControl/config/oauth';
import MenuComponent from '../../MenuComponent.vue'

interface ProviderView {
  name: OAuthProviderName
  icon: string
  color: string
  label: string
  enabled: boolean
  clientId?: string
  clientSecret?: string
}

const $route = useRoute()
const displayableNamespace = $route.params.currentNamespace === "_global" ? "" : $route.params.currentNamespace as string

const loading = ref(true)
const loadingError = ref('')

const providers = ref([
  {
    name: "google",
    color: "blue",
    icon: "mdi-google",
    label: "Google",
    enabled: true
  },
  {
    name: "facebook",
    color: "blue",
    icon: "mdi-facebook",
    label: "Facebook",
    enabled: true
  },
  {
    name: "github",
    color: "blue",
    icon: "mdi-github",
    label: "GitHub",
    enabled: true
  },
  {
    name: "gitlab",
    color: "blue",
    icon: "mdi-gitlab",
    label: "GitLab",
    enabled: true
  },
  {
    name: "twitter",
    color: "blue",
    icon: "mdi-twitter",
    label: "Twitter",
    enabled: true
  },
  {
    name: "microsoft",
    color: "blue",
    icon: "mdi-microsoft",
    label: "Microsoft",
    enabled: true
  },
  /*{
    name: "amazon",
    color: "blue",
    icon: "mdi-amazon",
    label: "Amazon",
    enabled: true
  },*/
  {
    name: "apple",
    color: "white",
    icon: "mdi-apple",
    label: "Apple",
    enabled: false
  },
  {
    name: "discord",
    color: "blue",
    icon: "mdi-discord",
    label: "Discord",
    enabled: true
  },
] as ProviderView[])

const selectedProviderName = ref('google' as OAuthProviderName)
const selectedProviderLabel = ref('Google')
const selectedProviderClientID = ref('')
const selectedProviderSecret = ref('')
const selectedProviderEnabled = ref(false)

async function loadData() {
    loading.value = true
    try {
        const response = await api.accessControl.config.oauth.getProvidersConfigs({ namespace: displayableNamespace })
        const configsMap = new Map<string, OAuthProviderConfig>()
        for (const config of response.providers) {
          configsMap.set(config.name, config)
        }
        for (const provider of providers.value) {
          const config = configsMap.get(provider.name)
          if (config !== undefined) {
            provider.enabled = config.enabled
            provider.clientId = config.clientID
            provider.clientSecret = config.clientSecret
          } else {
            provider.enabled = false
            provider.clientId = undefined
            provider.clientSecret = undefined
          }
        }
        selectProvider(providers.value[0])
    } catch (e) {
    } finally {
        loading.value = false
    }
}

async function updateProviderConfig() {
  try {
    const response = await api.accessControl.config.oauth.updateProviderConfig({
      namespace: displayableNamespace,
      name: selectedProviderName.value,
      authURL: "",
      clientID: selectedProviderClientID.value,
      clientSecret: selectedProviderSecret.value,
      enabled: selectedProviderEnabled.value,
      tokenURL: "",
      userApiURL: "",
    })
  } catch(e) {
    console.error(e)
  }

  await loadData()
}

function selectProvider(provider: ProviderView) {
  selectedProviderName.value = provider.name
  selectedProviderLabel.value = provider.label
  selectedProviderEnabled.value = provider.enabled
  selectedProviderClientID.value = provider.clientId || ""
  selectedProviderSecret.value = provider.clientSecret || ""
}

onMounted(async () => {
  await loadData()
})

</script>