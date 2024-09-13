import { post } from '$api/post';

export function postSelectedRegion(selectedRegion: number) {
	post('parameter/selected-region', { selectedRegion });
}
