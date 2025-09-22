<script lang="ts">
  import { createQuery } from "@tanstack/svelte-query";
  import type { ListDevicesRequest } from "#/gen/claw/v1/device_service_pb";
  import { getDeviceServiceClient } from "#/connectrpc";
  import type { M } from "#/types";
  import { fromQuery, toQuery } from "query-string-parser";
  import { watch, useDebounce } from "runed";
  import IconMonitorSmartphone from "@lucide/svelte/icons/monitor-smartphone";
  import IconCircleMinus from "@lucide/svelte/icons/circle-minus";

  // TODO: Fix query state. Use nuqs-svelte library.
  let queryState: Partial<M<ListDevicesRequest>> = $state(
    fromQuery(window.location.search) || {
      lastImage: {
        include: true,
      },
      countImages: true,
    },
  );

  const listDeviceQuery = createQuery({
    queryKey: ["devices", "list"],
    queryFn: async ({ signal }) => {
      console.log($state.snapshot(queryState));
      const client = await getDeviceServiceClient({ signal });
      return client.listDevices(queryState);
    },
    refetchOnWindowFocus() {
      return !!queryState.pagination?.prevToken;
    },
    refetchOnReconnect() {
      return !!queryState.pagination?.prevToken;
    },
  });

  const debounced = useDebounce(() => $listDeviceQuery.refetch());

  watch(
    () => queryState,
    (val, prev) => {
      if (!prev) return;
      const qs = toQuery(val);
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
      + Add
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
          <svg
            class="w-4 h-4"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M12 6v6m0 0v6m0-6h6m-6 0H6"
            ></path>
          </svg>
          Add Device
        </button>
      </div>
    {/if}
  {/if}
</div>
