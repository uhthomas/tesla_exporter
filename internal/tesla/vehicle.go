package tesla

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strconv"
)

type vehicleResponse struct {
	Response Vehicle `json:"response"`
}

type Vehicle struct {
	ID                     uint64              `json:"id"`
	UserID                 uint64              `json:"user_id"`
	VehicleID              uint64              `json:"vehicle_id"`
	VIN                    string              `json:"vin"`
	DisplayName            string              `json:"display_name"`
	OptionCodes            string              `json:"option_codes"`
	Color                  interface{}         `json:"color"`
	AccessType             string              `json:"access_type"`
	Tokens                 []string            `json:"tokens"`
	State                  string              `json:"state"`
	InService              bool                `json:"in_service"`
	IDS                    string              `json:"id_s"`
	CalendarEnabled        bool                `json:"calendar_enabled"`
	APIVersion             int                 `json:"api_version"`
	BackseatToken          interface{}         `json:"backseat_token"`
	BackseatTokenUpdatedAt interface{}         `json:"backseat_token_updated_at"`
	VehicleConfig          VehicleConfig       `json:"vehicle_config"`
	ChargeState            VehicleChargeState  `json:"charge_state"`
	ClimateState           VehicleClimateState `json:"climate_state"`
	DriveState             VehicleDriveState   `json:"drive_state"`
	GuiSettings            VehicleGuiSettings  `json:"gui_settings"`
	VehicleState           VehicleState        `json:"vehicle_state"`
}

type VehicleConfig struct {
	CanAcceptNavigationRequests bool        `json:"can_accept_navigation_requests"`
	CanActuateTrunks            bool        `json:"can_actuate_trunks"`
	CarSpecialType              string      `json:"car_special_type"`
	CarType                     string      `json:"car_type"`
	ChargePortType              string      `json:"charge_port_type"`
	DefaultChargeToMax          bool        `json:"default_charge_to_max"`
	EceRestrictions             bool        `json:"ece_restrictions"`
	EuVehicle                   bool        `json:"eu_vehicle"`
	ExteriorColor               string      `json:"exterior_color"`
	ExteriorTrim                string      `json:"exterior_trim"`
	HasAirSuspension            bool        `json:"has_air_suspension"`
	HasLudicrousMode            bool        `json:"has_ludicrous_mode"`
	KeyVersion                  int         `json:"key_version"`
	MotorizedChargePort         bool        `json:"motorized_charge_port"`
	Plg                         bool        `json:"plg"`
	RearSeatHeaters             int         `json:"rear_seat_heaters"`
	RearSeatType                interface{} `json:"rear_seat_type"`
	Rhd                         bool        `json:"rhd"`
	RoofColor                   string      `json:"roof_color"`
	SeatType                    interface{} `json:"seat_type"`
	SpoilerType                 string      `json:"spoiler_type"`
	SunRoofInstalled            interface{} `json:"sun_roof_installed"`
	ThirdRowSeats               string      `json:"third_row_seats"`
	Timestamp                   uint64      `json:"timestamp"`
	UseRangeBadging             bool        `json:"use_range_badging"`
	WheelType                   string      `json:"wheel_type"`
}

type VehicleChargeState struct {
	BatteryHeaterOn bool `json:"battery_heater_on"`
	// Always an integer, but we decode it as a float64 to convert to a
	// ratio later.
	BatteryLevel                float64     `json:"battery_level"`
	BatteryRange                float64     `json:"battery_range"`
	ChargeCurrentRequest        int         `json:"charge_current_request"`
	ChargeCurrentRequestMax     int         `json:"charge_current_request_max"`
	ChargeEnableRequest         bool        `json:"charge_enable_request"`
	ChargeEnergyAdded           float64     `json:"charge_energy_added"`
	ChargeLimitSoc              int         `json:"charge_limit_soc"`
	ChargeLimitSocMax           int         `json:"charge_limit_soc_max"`
	ChargeLimitSocMin           int         `json:"charge_limit_soc_min"`
	ChargeLimitSocStd           int         `json:"charge_limit_soc_std"`
	ChargeMilesAddedIdeal       float64     `json:"charge_miles_added_ideal"`
	ChargeMilesAddedRated       float64     `json:"charge_miles_added_rated"`
	ChargePortColdWeatherMode   bool        `json:"charge_port_cold_weather_mode"`
	ChargePortDoorOpen          bool        `json:"charge_port_door_open"`
	ChargePortLatch             string      `json:"charge_port_latch"`
	ChargeRate                  float64     `json:"charge_rate"`
	ChargeToMaxRange            bool        `json:"charge_to_max_range"`
	ChargerActualCurrent        int         `json:"charger_actual_current"`
	ChargerPhases               int         `json:"charger_phases"`
	ChargerPilotCurrent         int         `json:"charger_pilot_current"`
	ChargerPower                int         `json:"charger_power"`
	ChargerVoltage              int         `json:"charger_voltage"`
	ChargingState               string      `json:"charging_state"`
	ConnChargeCable             string      `json:"conn_charge_cable"`
	EstBatteryRange             float64     `json:"est_battery_range"`
	FastChargerBrand            string      `json:"fast_charger_brand"`
	FastChargerPresent          bool        `json:"fast_charger_present"`
	FastChargerType             string      `json:"fast_charger_type"`
	IdealBatteryRange           float64     `json:"ideal_battery_range"`
	ManagedChargingActive       bool        `json:"managed_charging_active"`
	ManagedChargingStartTime    interface{} `json:"managed_charging_start_time"`
	ManagedChargingUserCanceled bool        `json:"managed_charging_user_canceled"`
	MaxRangeChargeCounter       int         `json:"max_range_charge_counter"`
	MinutesToFullCharge         int         `json:"minutes_to_full_charge"`
	NotEnoughPowerToHeat        interface{} `json:"not_enough_power_to_heat"`
	ScheduledChargingPending    bool        `json:"scheduled_charging_pending"`
	ScheduledChargingStartTime  interface{} `json:"scheduled_charging_start_time"`
	TimeToFullCharge            float64     `json:"time_to_full_charge"`
	Timestamp                   uint64      `json:"timestamp"`
	TripCharging                bool        `json:"trip_charging"`
	// Always an integer, but we decode it as a float64 to convert to a
	// ratio later.
	UsableBatteryLevel      float64     `json:"usable_battery_level"`
	UserChargeEnableRequest interface{} `json:"user_charge_enable_request"`
}

type VehicleClimateState struct {
	BatteryHeater              bool        `json:"battery_heater"`
	BatteryHeaterNoPower       interface{} `json:"battery_heater_no_power"`
	ClimateKeeperMode          string      `json:"climate_keeper_mode"`
	DefrostMode                int         `json:"defrost_mode"`
	DriverTempSetting          float64     `json:"driver_temp_setting"`
	FanStatus                  int         `json:"fan_status"`
	InsideTemp                 float64     `json:"inside_temp"`
	IsAutoConditioningOn       bool        `json:"is_auto_conditioning_on"`
	IsClimateOn                bool        `json:"is_climate_on"`
	IsFrontDefrosterOn         bool        `json:"is_front_defroster_on"`
	IsPreconditioning          bool        `json:"is_preconditioning"`
	IsRearDefrosterOn          bool        `json:"is_rear_defroster_on"`
	LeftTempDirection          int         `json:"left_temp_direction"`
	MaxAvailTemp               float64     `json:"max_avail_temp"`
	MinAvailTemp               float64     `json:"min_avail_temp"`
	OutsideTemp                float64     `json:"outside_temp"`
	PassengerTempSetting       float64     `json:"passenger_temp_setting"`
	RemoteHeaterControlEnabled bool        `json:"remote_heater_control_enabled"`
	RightTempDirection         int         `json:"right_temp_direction"`
	SeatHeaterLeft             int         `json:"seat_heater_left"`
	SeatHeaterRearCenter       int         `json:"seat_heater_rear_center"`
	SeatHeaterRearLeft         int         `json:"seat_heater_rear_left"`
	SeatHeaterRearRight        int         `json:"seat_heater_rear_right"`
	SeatHeaterRight            int         `json:"seat_heater_right"`
	SideMirrorHeaters          bool        `json:"side_mirror_heaters"`
	Timestamp                  uint64      `json:"timestamp"`
	WiperBladeHeater           bool        `json:"wiper_blade_heater"`
}

type VehicleDriveState struct {
	GpsAsOf                 int         `json:"gps_as_of"`
	Heading                 int         `json:"heading"`
	Latitude                float64     `json:"latitude"`
	Longitude               float64     `json:"longitude"`
	NativeLatitude          float64     `json:"native_latitude"`
	NativeLocationSupported int         `json:"native_location_supported"`
	NativeLongitude         float64     `json:"native_longitude"`
	NativeType              string      `json:"native_type"`
	Power                   int         `json:"power"`
	ShiftState              interface{} `json:"shift_state"`
	Speed                   interface{} `json:"speed"`
	Timestamp               uint64      `json:"timestamp"`
}

type VehicleGuiSettings struct {
	Gui24HourTime       bool   `json:"gui_24_hour_time"`
	GuiChargeRateUnits  string `json:"gui_charge_rate_units"`
	GuiDistanceUnits    string `json:"gui_distance_units"`
	GuiRangeDisplay     string `json:"gui_range_display"`
	GuiTemperatureUnits string `json:"gui_temperature_units"`
	ShowRangeUnits      bool   `json:"show_range_units"`
	Timestamp           uint64 `json:"timestamp"`
}

type VehicleState struct {
	APIVersion               int                        `json:"api_version"`
	AutoparkStateV2          string                     `json:"autopark_state_v2"`
	AutoparkStyle            string                     `json:"autopark_style"`
	CalendarSupported        bool                       `json:"calendar_supported"`
	CarVersion               string                     `json:"car_version"`
	CenterDisplayState       int                        `json:"center_display_state"`
	Df                       int                        `json:"df"`
	Dr                       int                        `json:"dr"`
	FdWindow                 int                        `json:"fd_window"`
	FpWindow                 int                        `json:"fp_window"`
	Ft                       int                        `json:"ft"`
	IsUserPresent            bool                       `json:"is_user_present"`
	LastAutoparkError        string                     `json:"last_autopark_error"`
	Locked                   bool                       `json:"locked"`
	MediaState               VehicleMediaState          `json:"media_state"`
	NotificationsSupported   bool                       `json:"notifications_supported"`
	Odometer                 float64                    `json:"odometer"`
	ParsedCalendarSupported  bool                       `json:"parsed_calendar_supported"`
	Pf                       int                        `json:"pf"`
	Pr                       int                        `json:"pr"`
	RdWindow                 int                        `json:"rd_window"`
	RemoteStart              bool                       `json:"remote_start"`
	RemoteStartEnabled       bool                       `json:"remote_start_enabled"`
	RemoteStartSupported     bool                       `json:"remote_start_supported"`
	RpWindow                 int                        `json:"rp_window"`
	Rt                       int                        `json:"rt"`
	SentryMode               bool                       `json:"sentry_mode"`
	SentryModeAvailable      bool                       `json:"sentry_mode_available"`
	SmartSummonAvailable     bool                       `json:"smart_summon_available"`
	SoftwareUpdate           VehicleSoftwareUpdateState `json:"software_update"`
	SpeedLimitMode           VehicleSpeedLimitModeState `json:"speed_limit_mode"`
	SummonStandbyModeEnabled bool                       `json:"summon_standby_mode_enabled"`
	Timestamp                uint64                     `json:"timestamp"`
	ValetMode                bool                       `json:"valet_mode"`
	ValetPinNeeded           bool                       `json:"valet_pin_needed"`
	VehicleName              string                     `json:"vehicle_name"`
}

type VehicleMediaState struct {
	RemoteControlEnabled bool `json:"remote_control_enabled"`
}

type VehicleSoftwareUpdateState struct {
	DownloadPerc        int    `json:"download_perc"`
	ExpectedDurationSec int    `json:"expected_duration_sec"`
	InstallPerc         int    `json:"install_perc"`
	Status              string `json:"status"`
	Version             string `json:"version"`
}

type VehicleSpeedLimitModeState struct {
	Active          bool    `json:"active"`
	CurrentLimitMph float64 `json:"current_limit_mph"`
	MaxLimitMph     int     `json:"max_limit_mph"`
	MinLimitMph     int     `json:"min_limit_mph"`
	PinCodeSet      bool    `json:"pin_code_set"`
}

func (c *Client) Vehicle(ctx context.Context, id uint64) (*Vehicle, error) {
	u := *c.baseURL
	u.Path = path.Join(u.Path, "vehicles", strconv.FormatUint(id, 10), "vehicle_data")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("User-Agent", "tesla_exporter")

	res, err := c.c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d", res.StatusCode)
	}

	var out vehicleResponse
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("json decode: %w", err)
	}
	return &out.Response, nil
}
