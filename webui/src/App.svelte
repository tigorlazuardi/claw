<script lang="ts">
  import Sidebar from "./surfaces/Sidebar.svelte";
  import Routing from "./Routing.svelte";
  import { RouterView } from "@dvcol/svelte-simple-router/components";
</script>

<Routing>
  <main
    data-theme="nord"
    class="flex h-screen w-screen bg-base-100 text-base-content m-0 p-0"
  >
    <Sidebar>
      <RouterView>
        {#snippet loading()}
          {#await import("./components/RoutingTransition.svelte") then { default: RoutingTransition }}
            <RoutingTransition />
          {/await}
        {/snippet}
      </RouterView>
    </Sidebar>
  </main>
</Routing>
