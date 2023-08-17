<template>
    <q-card>
        <q-card-section>
            <p class="text-h6">Bind Balena device</p>
            <DeviceSelectorInput
                :namespace="namespaceInput"
                v-model="deviceInput"
            />
        </q-card-section>

        <q-card-actions>
            <q-btn
                outline
                label="Bind"
                color="dark"
                class="fit"
                :disable="binding || deviceInput === null"
                :loading="binding"
                @click="bind"
            />
        </q-card-actions>
    </q-card>
</template>

<script setup lang="ts">
import { ref } from 'vue';

import DeviceSelectorInput from '../../device/DeviceSelectorInput.vue';
import { Device } from 'src/boot/api/iot/device';
import { useQuasar } from 'quasar';
import { useI18n } from 'vue-i18n';
import api from 'src/boot/api';

const $q = useQuasar()
const $i18n = useI18n()

const props = defineProps<
    {
        balenaDeviceUUID: string;
    }
>()
const emits = defineEmits<{
    (e: 'onBinded', device: Device): void
}>()


const namespaceInput = ref("")
const deviceInput = ref(null as Device | null)

const binding = ref(false)

async function bind() {
    if (deviceInput.value === null) return
    binding.value = true

    const notif = $q.notify({
        message: $i18n.t('modules.iot.integration.balena.device.bindModal.'),
        type: 'ongoing',
    })

    try {
        console.log(deviceInput.value)

        await api.iot.integration.balena.device.bind({
            balenaDeviceUUID: props.balenaDeviceUUID,
            deviceNamespace: namespaceInput.value,
            deviceUuid: deviceInput.value.uuid
        })
        notif({
            message: $i18n.t('modules.iot.integration.balena.device.bindModal.bindSuccessNotify'),
            type: 'positive',
        })
        emits('onBinded', deviceInput.value)
    } catch (e) {
        notif({
            message: $i18n.t('modules.iot.integration.balena.device.bindModal.bindErrorNotify', { error: e }),
            type: 'negative',
        })
    } finally {
        binding.value = false
    }
}

</script>