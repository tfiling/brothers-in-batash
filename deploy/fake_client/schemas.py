import marshmallow as mw


class SoldierRoleSchema(mw.Schema):
    id = mw.fields.String(required=True)
    name = mw.fields.String(required=True)
    description = mw.fields.String(validate=mw.validate.Length(max=255))


class SoldierSchema(mw.Schema):
    id = mw.fields.String(required=True)
    first_name = mw.fields.String(required=True)
    middle_name = mw.fields.String()
    last_name = mw.fields.String(required=True)
    personal_number = mw.fields.String(required=True, validate=mw.validate.Regexp(r'^[0-9]{7}$'))
    position = mw.fields.Integer(required=True)
    roles = mw.fields.List(mw.fields.Nested(SoldierRoleSchema), validate=mw.validate.Length(min=1))


class ShiftSchema(mw.Schema):
    start_time_hour = mw.fields.Integer(required=True, validate=mw.validate.Range(min=0, max=23))
    start_time_minute = mw.fields.Integer(required=True, validate=mw.validate.Range(min=0, max=59))
    end_time_hour = mw.fields.Integer(required=True, validate=mw.validate.Range(min=0, max=23))
    end_time_minute = mw.fields.Integer(required=True, validate=mw.validate.Range(min=0, max=59))
    id = mw.fields.String(required=True)
    name = mw.fields.String(required=True)
    type = mw.fields.Integer(required=True, validate=mw.validate.Range(min=0))
    commander_soldier_id = mw.fields.String(required=True)
    additional_soldiers_ids = mw.fields.List(mw.fields.String())
    description = mw.fields.String(validate=mw.validate.Length(max=255))
    shift_template_id = mw.fields.String()


# TODO - apply additional adaptations into separated API body types
class TimeOfDaySchema(mw.Schema):
    hour = mw.fields.Integer(validate=mw.validate.Range(min=0, max=23))
    minute = mw.fields.Integer(validate=mw.validate.Range(min=0, max=59))


class ShiftTimeSchema(mw.Schema):
    start_time = mw.fields.Nested(TimeOfDaySchema, required=True)
    end_time = mw.fields.Nested(TimeOfDaySchema, required=True)


class PersonnelRequirementSchema(mw.Schema):
    soldier_role_to_count = mw.fields.Dict(keys=mw.fields.String(), values=mw.fields.Integer())


class ShiftTemplateSchema(mw.Schema):
    id = mw.fields.String(required=True)
    name = mw.fields.String(required=True)
    description = mw.fields.String(required=True)
    personnel_requirement = mw.fields.Nested(PersonnelRequirementSchema, required=True)
    days_of_occurrences = mw.fields.Dict(
        keys=mw.fields.Integer(validate=validate.Range(min=0, max=6)),
        values=mw.fields.List(mw.fields.Nested(ShiftTimeSchema)),
        required=True
    )
