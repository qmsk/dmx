export interface HeadType {
  Vendor: string;
  Model:  string;
  Mode:   string;
}
export interface HeadConfig {
  Type:     string;
  Universe: number;
  Address:  number;
}
export interface Head {
  ID:     string;
  Type:   HeadType;
  Config: HeadConfig;
}
