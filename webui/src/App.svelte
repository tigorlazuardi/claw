<script lang="ts">
  import Sidebar from "./surfaces/Sidebar.svelte";
  import Routing from "./Routing.svelte";
  import { RouterView } from "@dvcol/svelte-simple-router/components";
</script>

<Routing>
  <main class="flex h-screen w-screen">
    <Sidebar />
    <div class="flex-1">
      <RouterView>
        {#snippet loading()}
          {#await import("./components/RoutingTransition.svelte") then { default: RoutingTransition }}
            <RoutingTransition />
          {/await}
        {/snippet}
      </RouterView>
    </div>
  </main>
</Routing>
