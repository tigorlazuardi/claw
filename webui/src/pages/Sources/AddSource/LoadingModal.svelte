<script lang="ts">
  import DialogModal from "#/components/DialogModal.svelte";
</script>

<DialogModal open={true}>
  {#snippet title()}
    <span class="font-bold">Loading Data...</span>
  {/snippet}
  {#snippet description()}
    <div class="loading loading-ring text-primary w-[20rem]"></div>
  {/snippet}
</DialogModal>
