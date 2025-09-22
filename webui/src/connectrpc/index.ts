export type SourceServiceClient = Awaited<
  ReturnType<typeof getSourceServiceClient>
>;

export async function getSourceServiceClient(options?: RequestInit) {
  const { createClient } = await import("@connectrpc/connect");
  const { createConnectTransport } = await import("@connectrpc/connect-web");
  const { SourceService } = await import("#/gen/claw/v1/source_service_pb");
  const transport = createConnectTransport({
    baseUrl: import.meta.env.BASE_URL,
    fetch: (input, init) => {
      return fetch(input, {
        ...init,
        ...options,
        credentials: options?.credentials || "include",
      });
    },
  });
  return createClient(SourceService, transport);
}

export async function getDeviceServiceClient(options?: RequestInit) {
  const { createClient } = await import("@connectrpc/connect");
  const { createConnectTransport } = await import("@connectrpc/connect-web");
  const { DeviceService } = await import("#/gen/claw/v1/device_service_pb");
  const transport = createConnectTransport({
    baseUrl: import.meta.env.BASE_URL,
    fetch: (input, init) => {
      return fetch(input, {
        ...init,
        ...options,
        credentials: options?.credentials || "include",
      });
    },
  });
  return createClient(DeviceService, transport);
}
