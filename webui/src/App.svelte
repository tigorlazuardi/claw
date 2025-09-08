<script lang="ts">
  import Sidebar from "./surfaces/Sidebar.svelte";
  import Routing from "./Routing.svelte";
  import { RouterView } from "@dvcol/svelte-simple-router/components";
</script>

<Routing>
  <main
    data-theme="winter"
    class="flex h-screen w-screen bg-gray-900 text-gray-100 m-0 p-0"
  >
    <Sidebar>
      <RouterView />
    </Sidebar>
  </main>
</Routing>
