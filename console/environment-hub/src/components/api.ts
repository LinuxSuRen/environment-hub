function DefaultResponseProcess(response: any) {
    if (!response.ok) {
      switch (response.status) {
        case 401:
          throw new Error("Unauthenticated")
      }
      throw new Error(response.statusText)
    } else {
      return response.json()
    }
}

export function GetClusters(callback: (d: any) => void, errHandle?: (e: any) => void | null) {
    fetch('/v1/k3d/clusters')
    .then(DefaultResponseProcess)
    .then(callback).catch(errHandle)
}

export function GetCluster(name: string, callback: (d: any) => void, errHandle?: (e: any) => void | null) {
    fetch('/v1/k3d/clusters/' + name)
    .then(DefaultResponseProcess)
    .then(callback).catch(errHandle)
}

export function CreateClusters(data: any, callback: (d: any) => void, errHandle?: (e: any) => void | null) {
    const requestOptions = {
        method: 'POST',
        body: JSON.stringify(data)
    }

    fetch('/v1/k3d/clusters', requestOptions)
    .then(DefaultResponseProcess)
    .then(callback).catch(errHandle)
}

export function DeleteCluster(name: string, callback: (d: any) => void, errHandle?: (e: any) => void | null) {
  const requestOptions = {
    method: 'DELETE'
  }

  fetch('/v1/k3d/clusters/' + name, requestOptions)
  .then(DefaultResponseProcess)
  .then(callback).catch(errHandle)
}

export function GetKindClusters(callback: (d: any) => void, errHandle?: (e: any) => void | null) {
  fetch('/v1/kind/clusters')
  .then(DefaultResponseProcess)
  .then(callback).catch(errHandle)
}

export function CreateKindClusters(data: any, callback: (d: any) => void, errHandle?: (e: any) => void | null) {
  const requestOptions = {
      method: 'POST',
      body: JSON.stringify(data)
  }

  fetch('/v1/kind/clusters', requestOptions)
  .then(DefaultResponseProcess)
  .then(callback).catch(errHandle)
}

export function DeleteKindCluster(name: string, callback: (d: any) => void, errHandle?: (e: any) => void | null) {
const requestOptions = {
  method: 'DELETE'
}

fetch('/v1/kind/clusters/' + name, requestOptions)
.then(DefaultResponseProcess)
.then(callback).catch(errHandle)
}
