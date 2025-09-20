<script lang="ts">
  let showAddModal = $state(false);

  import { createQuery } from "@tanstack/svelte-query";
  import type {
    ListDevicesRequest
  } from "#/gen/claw/v1/device_service_pb";
</script>

<div class="p-[2rem] 2xl:max-w-[60vw] max-w-full m-auto">
  <div class="flex justify-between items-start mb-[2rem]">
    <div>
      <h1 class="text-4xl text-base-content">Sources</h1>
      <p class="text-lg font-light text-base-content/70">
        Configure and manage your image sources
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
</div>
