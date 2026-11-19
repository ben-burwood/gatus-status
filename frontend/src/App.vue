<template>
    <div class="min-h-screen w-full flex flex-col bg-base-200">
        <!-- HEADER -->
        <header class="bg-base-100 shadow-sm">
            <div class="max-w-3xl mx-auto px-5 py-5 flex items-center">
                <h1 class="text-2xl font-bold">{{ title }}</h1>
                <div class="flex-1"></div>
                <span v-if="updatedAt" class="text-sm text-base-content/60">Updated {{ updatedAt.toLocaleTimeString() }}</span>
            </div>
        </header>

        <!-- BODY -->
        <main class="flex-1 p-5">
            <div class="max-w-3xl mx-auto flex flex-col gap-6">
                <div v-if="error" role="alert" class="alert alert-warning">
                    <IconTriangleAlert class="size-5" />
                    <span>Unable to load the current status. Retrying…</span>
                </div>

                <template v-if="statuses.length">
                    <OverallStatus :statuses="statuses" />
                    <div class="card bg-base-100 shadow-sm">
                        <ul class="divide-y divide-base-200">
                            <li v-for="status in statuses" :key="status.group + '/' + status.name">
                                <EndpointRow :status="status" />
                            </li>
                        </ul>
                    </div>
                </template>
                <div v-else-if="loading" class="text-center mt-10">
                    <span class="loading loading-dots loading-xl"></span>
                </div>
            </div>
        </main>

        <!-- FOOTER -->
        <footer class="border-t border-base-300 bg-base-100">
            <div class="px-4 py-6 text-sm text-base-content/70 text-center">
                Powered by
                <a href="https://github.com/TwiN/gatus" target="_blank" rel="noopener" class="font-medium link link-hover">Gatus</a>
                <span class="mx-1">·</span>
                <a href="https://github.com/ben-burwood/gatus-status" target="_blank" rel="noopener" class="font-medium link link-hover">gatus-status</a>
            </div>
        </footer>
    </div>
</template>

<script setup lang="ts">
import EndpointRow from "@/components/EndpointRow.vue";
import OverallStatus from "@/components/OverallStatus.vue";
import { SERVER_URL } from "@/main";
import type { EndpointStatus } from "@/types/status";
import { onMounted, onUnmounted, ref } from "vue";
import IconTriangleAlert from "~icons/lucide/triangle-alert";

const REFRESH_INTERVAL_MS = 60 * 1000;

const title = ref("Status");
const statuses = ref<EndpointStatus[]>([]);
const loading = ref(true);
const error = ref(false);
const updatedAt = ref<Date | null>(null);
let refreshHandler: ReturnType<typeof setInterval> | null = null;

const loadConfig = async () => {
    try {
        const res = await fetch(`${SERVER_URL}/config`);
        const config: { title: string } = await res.json();
        title.value = config.title;
        document.title = config.title;
    } catch (e) {
        console.error("Failed to fetch config:", e);
    }
};

const loadStatuses = async () => {
    try {
        const res = await fetch(`${SERVER_URL}/statuses`);
        if (!res.ok) throw new Error(`HTTP ${res.status}`);
        statuses.value = await res.json();
        updatedAt.value = new Date();
        error.value = false;
    } catch (e) {
        console.error("Failed to fetch statuses:", e);
        error.value = true;
    } finally {
        loading.value = false;
    }
};

onMounted(() => {
    loadConfig();
    loadStatuses();
    refreshHandler = setInterval(loadStatuses, REFRESH_INTERVAL_MS);
});

onUnmounted(() => {
    if (refreshHandler) clearInterval(refreshHandler);
});
</script>
