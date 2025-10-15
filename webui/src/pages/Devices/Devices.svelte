<script lang="ts">
  import { createQuery } from "@tanstack/svelte-query";
  import { getDeviceServiceClient } from "#/connectrpc";
  import { getDevicePaginationSchema } from "#/connectrpc/pagination";
  import { fromQuery, toQuery } from "query-string-parser";
  import { watch, useDebounce } from "runed";
  import IconMonitorSmartphone from "@lucide/svelte/icons/monitor-smartphone";
  import IconCircleMinus from "@lucide/svelte/icons/circle-minus";
  import IconPlus from "@lucide/svelte/icons/plus";
  import AddDevice from "./AddDevice/AddDevice.svelte";

  import * as v from "valibot";
  import { nsfwState, parseNSFWState } from "#/store/searchQuery";

  function parseBool(s: string) {
    return s === "true";
  }

  function removeFalsyRecursive(obj: any): any {
    // Handle null or undefined input
    if (obj === null || obj === undefined) {
      return obj;
    }

    // Handle arrays
    if (Array.isArray(obj)) {
      return obj
        .map((item: any) => removeFalsyRecursive(item)) // Recursively process each item
        .filter((item) => Boolean(item)); // Remove falsy values
    }

    // Handle objects (but not Date, RegExp, etc.)
    if (typeof obj === "object" && obj.constructor === Object) {
      const result: any = {};

      for (const [key, value] of Object.entries(obj)) {
        const processedValue = removeFalsyRecursive(value);

        // Only add to result if the processed value is truthy
        if (Boolean(processedValue)) {
          result[key] = processedValue;
        }
      }

      return result;
    }

    // For primitive values, return as-is (filtering happens at parent level)
    return obj;
  }

  const querySchema = v.object({
    lastImage: v.fallback(
      v.object({
        include: v.fallback(v.pipe(v.string(), v.transform(parseBool)), true),
        nsfw: v.fallback(
          v.pipe(v.string(), v.transform(parseNSFWState)),
          nsfwState.current,
        ),
      }),
      {
        include: true,
        nsfw: nsfwState.current,
      },
    ),
    countImages: v.fallback(v.pipe(v.string(), v.transform(parseBool)), true),
    pagination: getDevicePaginationSchema(),
    search: v.fallback(v.string(), ""),
    sourceId: v.fallback(
      v.pipe(v.string(), v.transform(parseInt), v.number()),
      0,
    ),
  });

  let queryState = $state(
    v.parse(querySchema, fromQuery(window.location.search) || {}),
  );

  const listDeviceQuery = createQuery({
    queryKey: ["devices", "list"],
    queryFn: async ({ signal }) => {
      const client = await getDeviceServiceClient({ signal });
      return client.listDevices(queryState);
    },
    refetchOnWindowFocus() {
      return !!queryState.pagination.prevToken;
    },
    refetchOnReconnect() {
      return !!queryState.pagination.prevToken;
    },
  });

  const debounced = useDebounce(() => $listDeviceQuery.refetch());

  watch(
    () => queryState,
    (val, prev) => {
      console.log(toQuery(removeFalsyRecursive(val)));
      if (!prev) return;
      const qs = toQuery(removeFalsyRecursive(val));
      if (qs !== "") {
        window.history.replaceState(null, "", `?${qs}`);
      }
      debounced();
    },
  );

  let showAddModal = $state(false);
</script>

<div class="p-[2rem] 2xl:max-w-[60vw] max-w-full m-auto">
  <div class="flex justify-between items-start mb-[2rem]">
    <div>
      <h1 class="text-4xl text-base-content">Devices</h1>
      <p class="text-lg font-light text-base-content/70">
        Configure and manage your devices
      </p>
    </div>
    <button
      class="btn btn-lg btn-primary"
      onclick={() => (showAddModal = true)}
    >
      <IconPlus /> Add
    </button>
  </div>
  <div class="divider"></div>
  {#if $listDeviceQuery.isLoading}
    <div class="w-full flex justify-center items-center py-16">
      <div class="loading loading-ring w-[20vw] text-primary"></div>
    </div>
  {:else if $listDeviceQuery.isError}
    <div class="flex flex-col items-center justify-center">
      <div class="alert alert-error alert-soft mt-4">
        <IconCircleMinus />
        <span>{$listDeviceQuery.error.message}</span>
      </div>
    </div>
  {:else if $listDeviceQuery.data}
    {@const data = $listDeviceQuery.data}
    {#if data.items.length}{:else}
      <div class="flex flex-col items-center justify-center py-16 text-center">
        <div class="mb-6">
          <IconMonitorSmartphone
            class="mx-auto w-[6rem] h-[6rem] text-base-content"
          />
          <h3 class="text-xl font-semibold text-base-content mb-2">
            No Devices Have Been Added
          </h3>
          <p class="text-base-content/70 max-w-md">
            Get started by adding a device to manage your wallpapers and images.
          </p>
        </div>
        <button class="btn btn-primary" onclick={() => (showAddModal = true)}>
          <IconPlus />
          Add Device
        </button>
      </div>
    {/if}
  {/if}

  {#if showAddModal}
    <AddDevice open={showAddModal} />
  {/if}
</div>
