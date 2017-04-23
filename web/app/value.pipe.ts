import { Pipe, PipeTransform } from '@angular/core'

import { Value } from './types';


@Pipe({
  name: 'dmxValue',
})
export class ValuePipe implements PipeTransform {
  transform(value: Value, digits:number = 1): string {
    return (value * 100.0).toFixed(digits) + '%';
  }
}
