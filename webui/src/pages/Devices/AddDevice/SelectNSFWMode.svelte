<script lang="ts">
  import type { ClassValue } from "svelte/elements";
  import { NSFWMode } from "#/gen/claw/v1/nsfw_pb"; 
  import type { ComponentProps } from "@sjsf/form";

  type Props = ComponentProps["stringField"] & {
    class?: ClassValue;
  }


  let { value = $bindable(), class: className }: Props = $props();
</script>

<select class={['select', className]} bind:value>
  <option value={NSFWMode.NSFW_MODE_ALLOW}>Allow Both Non-NSFW and NSFW Content</option>
  <option value={NSFWMode.NSFW_MODE_DISALLOW}>Disallow NSFW Content</option>
  <option value={NSFWMode.NSFW_MODE_ONLY}>Only NSFW Content</option>
</select>
