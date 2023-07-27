<template>
    <div class="full-height full-width row justify-center items-center">
      <div class="row">
          <q-card square bordered class="q-pa-lg shadow-1">
              <q-card-section class="text-center">
                  <h3 class="q-ma-sm text-bold text-uppercase">{{ $t('modules.accessControl.iam.auth.certificate.register.header') }}</h3>
              </q-card-section>
          <q-card-section>
              <q-form class="q-gutter-md" ref="formRef">
              <q-input :disable="props.defaultNamespace != undefined" square filled v-model="namespace" counter maxlength="32" type="text" :label="$t('modules.accessControl.iam.auth.certificate.register.namespaceInput')" />
              <q-input :disable="props.defaultIdentityUUID != undefined" square filled clearable v-model="identityUUID" counter maxlength="32" type="text" :label="$t('modules.accessControl.iam.auth.certificate.register.identityInput')" />
              <q-input square filled clearable v-model="description" counter maxlength="128" type="text" :label="$t('modules.accessControl.iam.auth.certificate.register.descriptionInput')" />
              <div class="row">
                <div class="col-12">
                    {{ $t('modules.accessControl.iam.auth.certificate.register.publicKeyInfo') }}
                </div>
                <q-file
                    square filled
                    v-model="publicKey"
                    :label="$t('modules.accessControl.iam.auth.certificate.register.fileInput')"
                    class="col-6 q-mt-xs"
                    :rules="[val => !!val || $t('modules.accessControl.iam.auth.certificate.register.fileRequired')]"
                    clearable
                />
                <div class="col-1 q-mt-xs" />
                <q-btn
                    unelevated
                    color="dark"
                    :label="$t('modules.accessControl.iam.auth.certificate.register.generateKeyPairButton')"
                    class="col-5 q-mt-xs q-mb-md"
                    icon="add"
                    @click="generateKey"
                    :disable="publicKey != null"
                />
              </div>
              </q-form>
          </q-card-section>
          <q-card-actions class="q-px-md">
              <q-btn unelevated color="dark" size="lg" class="full-width" :label="$t('modules.accessControl.iam.auth.certificate.register.registerAndGenerateButton')" :loading="loading" :disabled="loading" @click="registerKeyAndGetCertificate" />
          </q-card-actions>
          <q-card-section class="text-center q-pa-none">
              <p class="text-grey-6">{{ $t('modules.accessControl.iam.auth.certificate.register.registerHint') }}</p>
          </q-card-section>
          </q-card>
      </div>
  </div>
</template>

<script lang="ts" setup>
import { AxiosError } from 'axios';
import { useQuasar } from 'quasar';
import api from 'src/boot/api';
import { Certificate } from 'src/boot/api/accessControl/auth/certificate';
import { onMounted, ref } from 'vue';
import { useI18n } from 'vue-i18n';

const $q = useQuasar()
const $i18n = useI18n()

const props = defineProps<{
    defaultNamespace?: string,
    defaultIdentityUUID?: string
}>()

const emit = defineEmits<{
  (e: 'canceled'): void,
  (e: 'registered', certificate: Certificate): void
}>()

const namespace = ref('')
const identityUUID = ref('')
const description = ref('')
const publicKey = ref(null as File | null)

const loading = ref(false)

const formRef = ref()

function download(filename: string, text: string) {
  var element = document.createElement('a');
  element.setAttribute('href', 'data:text/plain;charset=utf-8,' + encodeURIComponent(text));
  element.setAttribute('download', filename);

  element.style.display = 'none';
  document.body.appendChild(element);

  element.click();

  document.body.removeChild(element);
}

async function generateKey() {
    try {
    const c = window.crypto.subtle

    const keyPair = await c.generateKey({
        name: "RSA-OAEP",
        modulusLength: 2048,
        publicExponent: new Uint8Array([0x01, 0x00, 0x01]),
        hash: { name: "SHA-256" }
    }, true, ["encrypt", "decrypt"])


    function convertToPEM(key: ArrayBuffer, prefix: string, sufix: string) {
        var byteArray = new Uint8Array(key);
        var byteString = '';
        for(var i=0; i < byteArray.byteLength; i++) {
            byteString += String.fromCharCode(byteArray[i]);
        }
        var b64PrivateKey = window.btoa(byteString);

        var finalString = '';
        while(b64PrivateKey.length > 0) {
            finalString += b64PrivateKey.substring(0, 64) + '\n';
            b64PrivateKey = b64PrivateKey.substring(64);
        }

        return prefix + "\n" + finalString + sufix;
    }

    const rawPrivateKey = await c.exportKey(
        "pkcs8",
        keyPair.privateKey
    )
    const privateKeyPEM = convertToPEM(rawPrivateKey, "-----BEGIN PRIVATE KEY-----", "-----END PRIVATE KEY-----")
    download("private.pem", privateKeyPEM)

    const rawPublicKey = await c.exportKey(
        "spki",
        keyPair.publicKey
    )
    const publicKeyPEM = convertToPEM(rawPublicKey, "-----BEGIN PUBLIC KEY-----", "-----END PUBLIC KEY-----")
    publicKey.value = new File([publicKeyPEM], "public.pem", { type: 'application/x-pem-file' })
    } catch(err) {
        console.log("Err")
        console.error(err)
    }
}

async function registerKeyAndGetCertificate() {
    try {
        await formRef.value.validate()
    } catch {   
        return
    }

    const notif = $q.notify({
      type: 'ongoing',
      message: $i18n.t('modules.accessControl.iam.auth.certificate.register.registerOperationNotify')
    })
    loading.value = true
    try {
        const publicKeyText = await publicKey.value?.text()
        if (publicKeyText == undefined) {
            throw "Faile to read public key file"
        }

        const response = await api.accessControl.auth.certificate.registerAndGenerate({
            description: description.value,
            identityUUID: identityUUID.value,
            namespace: namespace.value,
            publicKey: publicKeyText
        })
        const certificate = response.raw
        download("client.crt", certificate)


        emit('registered', response.certificate)
        notif({
            type: 'positive',
            message: $i18n.t('modules.accessControl.iam.auth.certificate.register.registerSuccessNotify')
        })
    } catch (error) {
        const axiosErr = error as AxiosError<{message: string}>
        const errorMesage = axiosErr.response?.data.message
        
        const errorString = errorMesage ? errorMesage : String(error)

        notif({
            type: 'negative',
            message: $i18n.t('modules.accessControl.iam.auth.certificate.register.registerFailNotify', { error: errorString }),
            timeout: 5000
        })
    } finally {
        loading.value = false
    }
}

function applyDefaults() {
    if (props.defaultNamespace != undefined) {
        namespace.value = props.defaultNamespace
    }
    if (props.defaultIdentityUUID != undefined) {
        identityUUID.value = props.defaultIdentityUUID
    }
}

onMounted(() => applyDefaults())

</script>