<template>
    <q-page>
        <q-spinner></q-spinner>
    </q-page>
</template>

<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router'
import { LocalStorage, useQuasar } from 'quasar'
import { useI18n } from 'vue-i18n'
import api from 'src/boot/api';
import { onMounted } from 'vue';
import { OAuthProviderName } from 'src/boot/api/me/auth';

const $q = useQuasar()
const $route = useRoute()
const $router = useRouter()
const $i18n = useI18n()

// Fix vue-router bug
function findGetParameter(parameterName: string): string {
    var result = null,
        tmp = [];
    location.search
        .substr(1)
        .split("&")
        .forEach(function (item) {
          tmp = item.split("=");
          if (tmp[0] === parameterName) result = decodeURIComponent(tmp[1]);
        });
    return result as unknown as string;
}

async function handleAuth() {
    const provider = $route.params.provider as OAuthProviderName
    const code = findGetParameter('code')
    const state = findGetParameter('state')

    const notif = $q.notify({
        type: 'ongoing',
        message: $i18n.t('modules.me.auth.oauth.finalize.finalizeOperationNotify'),
    })

    const savedState = LocalStorage.getItem('me_auth_oauth_state') as string
    if (savedState !== state) {
        notif({
            type: 'negative',
            message: $i18n.t('modules.me.auth.oauth.finalize.finalizeInvalidStateNotify'),
            timeout: 5000
        })
        await $router.push({ name: 'me_auth' })
    }


    try {
        const response = await api.me.auth.finalizeOAuthRegistration({
            provider,
            code,
        })
        if (response.status == 'OK') {
            notif({
                type: 'positive',
                message: $i18n.t('modules.me.auth.oauth.finalize.finalizeSuccessNotify'),
                timeout: 5000
            }) 
        } else {
            notif({
                type: 'negative',
                message: $i18n.t('modules.me.auth.oauth.finalize.finalizeFailNotify', { error: response.status }),
                timeout: 5000
            })
        }
    } catch (e) {
        notif({
            type: 'negative',
            message: $i18n.t('modules.me.auth.oauth.finalize.finalizeFailNotify', { error: e }),
            timeout: 5000
        })
    } finally {
        await $router.push({ name: 'me_auth' })
    }
}

onMounted(async () => {
    await handleAuth()
})

</script>