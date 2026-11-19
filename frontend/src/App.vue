<template>
    <div class="min-h-screen w-full flex flex-col bg-zinc-50 text-zinc-900 dark:bg-zinc-950 dark:text-zinc-100">
        <main class="flex-1 w-full max-w-3xl mx-auto px-5 py-12 flex flex-col gap-6">
            <div class="text-center">
                <h1 class="text-3xl font-bold">{{ title }}</h1>
                <p v-if="updatedAt" class="mt-1 text-sm text-zinc-500">Updated {{ updatedAt.toLocaleTimeString() }}</p>
            </div>

            <div
                v-if="error"
                role="alert"
                class="flex items-center gap-3 rounded-lg px-4 py-3 bg-amber-100 text-amber-900 dark:bg-amber-950 dark:text-amber-200"
            >
                <IconTriangleAlert class="size-5 shrink-0" />
                <span>Unable to load the current status. Retrying…</span>
            </div>

            <template v-if="statuses.length">
                <OverallStatus :statuses="statuses" />
                <ul class="rounded-lg bg-white shadow-sm ring-1 ring-zinc-200 divide-y divide-zinc-200 dark:bg-zinc-900 dark:ring-zinc-800 dark:divide-zinc-800">
                    <li v-for="status in statuses" :key="status.group + '/' + status.name">
                        <EndpointRow :status="status" />
                    </li>
                </ul>
            </template>
            <div v-else-if="loading" class="flex justify-center mt-10">
                <IconLoaderCircle class="size-8 animate-spin text-zinc-400" />
            </div>
        </main>

        <footer class="py-6 flex items-center justify-center gap-3 text-sm text-zinc-500">
            <span>
                Powered by
                <a href="https://github.com/TwiN/gatus" target="_blank" rel="noopener" class="font-medium hover:underline">Gatus</a>
                <span class="mx-1">·</span>
                <a href="https://github.com/ben-burwood/gatus-status" target="_blank" rel="noopener" class="font-medium hover:underline">gatus-status</a>
            </span>
            <ThemeToggle />
        </footer>
    </div>
</template>

<script setup lang="ts">
import EndpointRow from "@/components/EndpointRow.vue";
import OverallStatus from "@/components/OverallStatus.vue";
import ThemeToggle from "@/components/ThemeToggle.vue";
import { SERVER_URL } from "@/main";
import type { EndpointStatus } from "@/types/status";
import { onMounted, onUnmounted, ref } from "vue";
import IconLoaderCircle from "~icons/lucide/loader-circle";
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
